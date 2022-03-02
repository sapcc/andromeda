/*
 *   Copyright 2021 SAP SE
 *
 *   Licensed under the Apache License, Version 2.0 (the "License");
 *   you may not use this file except in compliance with the License.
 *   You may obtain a copy of the License at
 *
 *       http://www.apache.org/licenses/LICENSE-2.0
 *
 *   Unless required by applicable law or agreed to in writing, software
 *   distributed under the License is distributed on an "AS IS" BASIS,
 *   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *   See the License for the specific language governing permissions and
 *   limitations under the License.
 */

package cli

type outputFormatters struct {
	Format     string   `short:"f" long:"format" description:"The output format, defaults to table" choice:"table" choice:"csv" choice:"markdown" choice:"html" choice:"value" default:"table"`
	Columns    []string `short:"c" long:"column" description:"specify the column(s) to include, can be repeated to show multiple columns"`
	SortColumn []string `long:"sort-column" description:"specify the column(s) to sort the data (columns specified first have a priority, non-existing columns are ignored), can be repeated"`
}
