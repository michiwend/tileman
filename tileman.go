/*
 * tileman - download a sequence of weather radar images from kachelmannwetter.com
 * Copyright © 2015 Michael Wendland
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
)

var regions = map[string]int{
	"germany": 2,
	"bavaria": 38,
}

var baseUrl = "http://kachelmannwetter.com/images/data/cache/"

func genSequence(start, end time.Time, region, resolution int) []string {

	start = start.Round(time.Minute * time.Duration(resolution))
	end = end.Round(time.Minute * time.Duration(resolution)) //FIXME set limit (now-limit)

	steps := int(end.Sub(start).Minutes()) / resolution

	t := start
	out := make([]string, steps+1)

	for i := 0; i <= steps; i++ {

		out[i] = fmt.Sprintf("download_px250_%s_%d_%s.png", t.Format("2006_01_02"), region, t.Format("1504"))
		t = t.Add(time.Duration(resolution) * time.Minute)

	}

	return out
}

func downloadSequence(start, end time.Time, region, resolution int, dir string) error {

	err := os.Mkdir(dir, 0775)
	if err != nil {
		return err
	}

	for i, img := range genSequence(start, end, region, resolution) {

		log.WithField("Name", img).Info("Downloading image...")

		resp, err := http.Get(baseUrl + img)
		if err != nil {
			return err
		}

		if resp.StatusCode != http.StatusOK {
			log.Error("Bad http response: " + resp.Status + " (" + img + ")")
		}

		buf, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		f, err := os.Create(path.Join(dir, strconv.Itoa(i)+"_"+img))
		if err != nil {
			return err
		}

		f.Write(buf)

		f.Close()
	}
	return nil
}

func main() {

	var startDateString = flag.String("start-date", time.Now().Format("2006-01-02"), `Start date in the form "2006-01-20", default is today.`)
	var endDateString = flag.String("end-date", time.Now().Format("2006-01-02"), `End date in the form "2006-01-20", default is today.`)

	var startTimeString = flag.String("start-time", time.Now().Add(time.Hour*-2).Format("15:04"), `Start time in the form "15:04", default is 2 hours ago.`)
	var endTimeString = flag.String("end-time", time.Now().Add(time.Minute*-15).Format("15:04"), `Start time in the form "15:04", default is 15 minutes ago.`)

	var outputDir = flag.String("dir", "./tileman_out", "Directory for saving the results.")

	var regionString = flag.String("region", "germany", "Which region map to use?")

	var resolution = flag.Int("res", 5, "Time resolution. Use a multiple of 5, minimum 5!")

	flag.Parse()

	if *resolution < 5 || *resolution%5 != 0 {
		log.Fatal("Invalid resolution value: " + strconv.Itoa(*resolution))
	}

	region, ok := regions[*regionString]
	if !ok {
		log.Fatal("Region not (yet) defined: " + *regionString)
	}

	start, err := time.Parse("2006-01-0215:04", *startDateString+*startTimeString)
	if err != nil {
		log.Fatal(err)
	}
	end, err := time.Parse("2006-01-0215:04", *endDateString+*endTimeString)
	if err != nil {
		log.Fatal(err)
	}

	err = downloadSequence(start, end, region, *resolution, *outputDir)
	if err != nil {
		log.Fatal(err)
	}
}
