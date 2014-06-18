/* Code based on project https://github.com/heroku/force
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"github.com/interactiveintelligence/icws_golib"
)

func DisplayConfigRecords(records []icws_golib.ConfigRecord) {
	if len(records) > 0 {
		fmt.Print(RenderConfigRecords(records))
	}
}

func recordColumns(records []icws_golib.ConfigRecord) (columns []string) {
	for _, record := range records {
		for key, _ := range record {
			found := false
			for _, column := range columns {
				if column == key {
					found = true
					break
				}
			}
			if !found && key != "configurationId.displayName" && key != "configurationId.uri" {
				columns = append(columns, key)
			}
		}
	}
	return
}

func coerceConfigRecords(uncoerced []map[string]interface{}) (records []icws_golib.ConfigRecord) {
	records = make([]icws_golib.ConfigRecord, len(uncoerced))
	for i, record := range uncoerced {
		records[i] = icws_golib.ConfigRecord(record)
	}
	return
}

func columnLengths(records []icws_golib.ConfigRecord, prefix string) (lengths map[string]int) {
	lengths = make(map[string]int)

	columns := recordColumns(records)
	for _, column := range columns {
		lengths[fmt.Sprintf("%s.%s", prefix, column)] = len(column) + 2
	}

	for _, record := range records {
		for column, value := range record {
			key := fmt.Sprintf("%s.%s", prefix, column)
			length := 0
			switch value := value.(type) {
			case []icws_golib.ConfigRecord:
				lens := columnLengths(value, key)
				for k, l := range lens {
					length += l
					if l > lengths[k] {
						lengths[k] = l
					}
				}
				length += len(lens) - 1
			case []interface{}:
				var buffer bytes.Buffer
				for _, element := range value {
					switch i := element.(type) {
					case map[string]interface{}:
						buffer.WriteString(fmt.Sprintf("%v,", i["id"]))
					default:
						buffer.WriteString(fmt.Sprintf("%v,", i))
					}
				}

				length += len(strings.Trim(buffer.String(), ","))
			default:
				if value == nil {
					length = len(" (null) ")
				} else {
					length = len(fmt.Sprintf(" %v ", value))
				}
			}
			if length > lengths[key] {
				lengths[key] = length
			}
		}
	}
	return
}

func recordHeader(columns []string, lengths map[string]int, prefix string) (out string) {
	headers := make([]string, len(columns))
	for i, column := range columns {
		key := fmt.Sprintf("%s.%s", prefix, column)
		headers[i] = fmt.Sprintf(fmt.Sprintf(" %%-%ds ", lengths[key]-2), column)
	}
	out = strings.Join(headers, "|")
	return
}

func recordSeparator(columns []string, lengths map[string]int, prefix string) (out string) {
	separators := make([]string, len(columns))
	for i, column := range columns {
		key := fmt.Sprintf("%s.%s", prefix, column)
		separators[i] = strings.Repeat("-", lengths[key])
	}
	out = strings.Join(separators, "+")
	return
}

func recordRow(record icws_golib.ConfigRecord, columns []string, lengths map[string]int, prefix string) (out string) {
	values := make([]string, len(columns))
	for i, column := range columns {
		value := record[column]
		switch value := value.(type) {
		case []icws_golib.ConfigRecord:
			values[i] = strings.TrimSuffix(renderConfigRecords(value, fmt.Sprintf("%s.%s", prefix, column), lengths), "\n")
		case []interface{}:
			var buffer bytes.Buffer
			for _, element := range value {
				switch i := element.(type) {
				case map[string]interface{}:
					buffer.WriteString(fmt.Sprintf("%v,", i["id"]))
				default:
					buffer.WriteString(fmt.Sprintf("%v,", i))
				}
			}

			values[i] = fmt.Sprintf(fmt.Sprintf(" %%-%dv ", lengths[column]-2), strings.Trim(buffer.String(), ","))
		default:
			if value == nil {
				values[i] = fmt.Sprintf(fmt.Sprintf(" %%-%ds ", lengths[column]-2), "(null)")
			} else {
				values[i] = fmt.Sprintf(fmt.Sprintf(" %%-%dv ", lengths[column]-2), value)
			}
		}
	}
	maxrows := 1
	for _, value := range values {
		rows := len(strings.Split(value, "\n"))
		if rows > maxrows {
			maxrows = rows
		}
	}
	rows := make([]string, maxrows)
	for i := 0; i < maxrows; i++ {
		rowvalues := make([]string, len(columns))
		for j, column := range columns {
			key := fmt.Sprintf("%s.%s", prefix, column)
			parts := strings.Split(values[j], "\n")
			if i < len(parts) {
				rowvalues[j] = fmt.Sprintf(fmt.Sprintf("%%-%ds", lengths[key]), parts[i])
			} else {
				rowvalues[j] = strings.Repeat(" ", lengths[key])
			}
		}
		rows[i] = strings.Join(rowvalues, "|")
	}
	out = strings.Join(rows, "\n")
	return
}

func flattenConfigRecord(record icws_golib.ConfigRecord) (flattened icws_golib.ConfigRecord) {
	flattened = make(icws_golib.ConfigRecord)
	for key, value := range record {
		if key == "attributes" {
			continue
		}
		switch value := value.(type) {
		case map[string]interface{}:
			if value["records"] != nil {
				unflattened := value["records"].([]interface{})
				subflattened := make([]icws_golib.ConfigRecord, len(unflattened))
				for i, record := range unflattened {
					subflattened[i] = (map[string]interface{})(flattenConfigRecord(icws_golib.ConfigRecord(record.(map[string]interface{}))))
				}
				flattened[key] = subflattened
			} else {
				for k, v := range flattenConfigRecord(value) {
					flattened[fmt.Sprintf("%s.%s", key, k)] = v
				}
			}
		default:
			flattened[key] = value
		}
	}
	return
}

func recordsHaveSubRows(records []icws_golib.ConfigRecord) bool {
	for _, record := range records {
		for _, value := range record {
			switch value := value.(type) {
			case []interface{}:
				if len(value) > 0 {
					return true

				}
			}
		}
	}
	return false
}

func renderConfigRecords(records []icws_golib.ConfigRecord, prefix string, lengths map[string]int) string {
	var out bytes.Buffer

	columns := recordColumns(records)

	out.WriteString(recordHeader(columns, lengths, prefix) + "\n")
	out.WriteString(recordSeparator(columns, lengths, prefix) + "\n")

	for _, record := range records {
		out.WriteString(recordRow(record, columns, lengths, prefix) + "\n")
		if recordsHaveSubRows(records) {

			out.WriteString(recordSeparator(columns, lengths, prefix) + "\n")
		}
	}

	return out.String()
}

func RenderConfigRecords(records []icws_golib.ConfigRecord) string {
	flattened := make([]icws_golib.ConfigRecord, len(records))
	for i, record := range records {
		flattened[i] = flattenConfigRecord(record)
	}
	lengths := columnLengths(flattened, "")
	return renderConfigRecords(flattened, "", lengths)
}

func DisplayConfigRecord(record icws_golib.ConfigRecord) {
	DisplayInterfaceMap(record, 0)
}

func DisplayInterfaceMap(object map[string]interface{}, indent int) {
	keys := make([]string, len(object))
	i := 0
	for key, _ := range object {
		keys[i] = key
		i++
	}
	sort.Strings(keys)
	for _, key := range keys {
		printIndent(indent)
		fmt.Printf("%s: ", key)
		switch v := object[key].(type) {
		case map[string]interface{}:
			fmt.Printf("\n")
			DisplayInterfaceMap(v, indent+1)
		case []interface{}:
			fmt.Printf("\n")
			for _, element := range v {
				switch i := element.(type) {
				case map[string]interface{}:
					DisplayInterfaceMap(i, indent+1)
				default:
					printIndent(indent + 1)
					fmt.Printf("%v\n", i)
				}
			}
		default:
			fmt.Printf("%v\n", v)
		}
	}
}

func DisplayInterfaceMapKeys(object map[string]interface{}) {
	keys := make([]string, len(object))
	i := 0
	for key, _ := range object {
		keys[i] = key
		i++
	}
	sort.Strings(keys)
	for _, key := range keys {
		fmt.Printf("%s\n", key)
	}
}

func DisplayList(list []string) {
	sort.Strings(list)
	for _, key := range list {
		fmt.Printf("%s\n", key)
	}
}

func printIndent(indent int) {
	for i := 0; i < indent; i++ {
		fmt.Printf("  ")
	}
}

func StringSliceToInterfaceSlice(s []string) (i []interface{}) {
	for _, str := range s {
		i = append(i, interface{}(str))
	}
	return
}
