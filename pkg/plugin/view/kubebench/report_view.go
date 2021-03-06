package kubebench

import (
	"fmt"
	"strconv"

	"github.com/aquasecurity/starboard-octant-plugin/pkg/plugin/view"

	starboard "github.com/aquasecurity/starboard/pkg/apis/aquasecurity/v1alpha1"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

// NewReport creates a new view component for displaying the specified CISKubeBenchReport.
func NewReport(benchmark *starboard.CISKubeBenchReport) (flexLayout component.FlexLayout) {
	flexLayout = *component.NewFlexLayout("CIS Kubernetes Benchmark")

	if benchmark == nil {
		flexLayout.AddSections(component.FlexLayoutSection{
			{
				Width: component.WidthFull,
				View: component.NewMarkdownText("This report is not available.\n" +
					"> Note that [kube-bench] reports are represented by instances of the `ciskubebenchreports.aquasecurity.github.io` resource.\n" +
					"> You can create such a report by running [kube-bench] with [Starboard CLI][starboard-cli]:\n" +
					"> ```\n" +
					"> $ starboard kube-bench\n" +
					"> ```\n" +
					"\n" +
					"[kube-bench]: https://github.com/aquasecurity/kube-bench\n" +
					"[starboard-cli]: https://github.com/aquasecurity/starboard#starboard-cli"),
			},
		})
		return
	}

	uiSections := make([]component.FlexLayoutItem, len(benchmark.Report.Sections))

	for i, section := range benchmark.Report.Sections {
		uiSections[i] = component.FlexLayoutItem{
			Width: component.WidthFull,
			View:  createTableForSection(section),
		}
	}

	uiSections = append([]component.FlexLayoutItem{
		{
			Width: component.WidthThird,
			View:  view.NewReportSummary(benchmark.CreationTimestamp.Time),
		},
		{
			Width: component.WidthThird,
			View:  view.NewScannerSummary(benchmark.Report.Scanner),
		},
		{
			Width: component.WidthThird,
			View:  NewCISKubeBenchReportSummary(benchmark),
		},
	}, uiSections...)

	flexLayout.AddSections(uiSections)
	return
}

func createTableForSection(section starboard.CISKubeBenchSection) component.Component {
	table := component.NewTableWithRows(
		fmt.Sprintf("%s %s", section.ID, section.Text), "There are no results!",
		component.NewTableCols("Status", "Number", "Description", "Scored"),
		[]component.TableRow{})

	for _, test := range section.Tests {
		for _, result := range test.Results {

			tr := component.TableRow{
				"Status":      component.NewText(result.Status),
				"Number":      component.NewText(result.TestNumber),
				"Description": component.NewText(result.TestDesc),
				"Scored":      component.NewText(strconv.FormatBool(result.Scored)),
			}
			table.Add(tr)
		}
	}

	return table
}

// TODO Implement summary counting
func NewCISKubeBenchReportSummary(report *starboard.CISKubeBenchReport) (summary *component.Summary) {
	totalPass := 0
	totalInfo := 0
	totalWarn := 0
	totalFail := 0

	for _, section := range report.Report.Sections {
		totalPass += section.TotalPass
		totalInfo += section.TotalInfo
		totalWarn += section.TotalWarn
		totalFail += section.TotalFail
	}

	summary = component.NewSummary("Summary")

	summary.Add([]component.SummarySection{
		{Header: "PASS ", Content: component.NewText(strconv.Itoa(totalPass))},
		{Header: "INFO", Content: component.NewText(strconv.Itoa(totalInfo))},
		{Header: "WARN ", Content: component.NewText(strconv.Itoa(totalWarn))},
		{Header: "FAIL ", Content: component.NewText(strconv.Itoa(totalFail))},
	}...)
	return
}
