package system

// import (
// 	"fmt"
// 	"io"
// 	"io/ioutil"
// 	"text/template"
// 	"time"
//
// 	"context"
//
// 	"github.com/dnephin/cobra"
// 	"github.com/storageos/go-cli/api/types"
// 	"github.com/storageos/go-cli/cli"
// 	"github.com/storageos/go-cli/cli/command"
// 	"github.com/storageos/go-cli/pkg/jsonlog"
// 	"github.com/storageos/go-cli/pkg/templates"
// )
//
// type eventsOptions struct {
// 	since  string
// 	until  string
// 	format string
// }
//
// // NewEventsCommand creates a new cobra.Command for `docker events`
// func NewEventsCommand(storageosCli *command.StorageOSCli) *cobra.Command {
// 	opt := eventsOptions{}
//
// 	cmd := &cobra.Command{
// 		Use:   "events [OPTIONS]",
// 		Short: "Get real time events from the server",
// 		Args:  cli.NoArgs,
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			return runEvents(storageosCli, &opt)
// 		},
// 	}
//
// 	flags := cmd.Flags()
// 	flags.StringVar(&opt.since, "since", "", "Show all events created since timestamp")
// 	flags.StringVar(&opt.until, "until", "", "Stream events until this timestamp")
// 	flags.StringVar(&opt.format, "format", "", "Format the output using the given Go template")
//
// 	return cmd
// }
//
// func runEvents(storageosCli *command.StorageOSCli, opt *eventsOptions) error {
// 	tmpl, err := makeTemplate(opt.format)
// 	if err != nil {
// 		return cli.StatusError{
// 			StatusCode: 64,
// 			Status:     "Error parsing format: " + err.Error()}
// 	}
// 	// options := types.EventsOptions{
// 	// 	Since:   opts.since,
// 	// 	Until:   opts.until,
// 	// 	Filters: opts.filter.Value(),
// 	// }
//
// 	ctx, cancel := context.WithCancel(context.Background())
// 	events, errs := storageosCli.Client().Events(ctx, types.ListOptions{})
// 	defer cancel()
//
// 	out := storageosCli.Out()
//
// 	for {
// 		select {
// 		case event := <-events:
// 			if err := handleEvent(out, event, tmpl); err != nil {
// 				return err
// 			}
// 		case err := <-errs:
// 			if err == io.EOF {
// 				return nil
// 			}
// 			return err
// 		}
// 	}
// }
//
// func handleEvent(out io.Writer, event types.Request, tmpl *template.Template) error {
// 	if tmpl == nil {
// 		return prettyPrintEvent(out, event)
// 	}
// 	return formatEvent(out, event, tmpl)
// }
//
// func makeTemplate(format string) (*template.Template, error) {
// 	if format == "" {
// 		return nil, nil
// 	}
// 	tmpl, err := templates.Parse(format)
// 	if err != nil {
// 		return tmpl, err
// 	}
// 	// we execute the template for an empty message, so as to validate
// 	// a bad template like "{{.badFieldString}}"
// 	return tmpl, tmpl.Execute(ioutil.Discard, &types.Event{})
// }
//
// // prettyPrintEvent prints all types of event information.
// // Each output includes the event type, actor id, name and action.
// // Actor attributes are printed at the end if the actor has any.
// func prettyPrintEvent(out io.Writer, event types.Request) error {
// 	if event.Timestamp != 0 {
// 		fmt.Fprintf(out, "%s ", time.Unix(0, event.Timestamp).Format(jsonlog.RFC3339NanoFixed))
// 	}
//
// 	fmt.Fprintf(out, "%s %s %s", event.EventType, event.Action, event.Message)
//
// 	// if len(event.Actor.Attributes) > 0 {
// 	// 	var attrs []string
// 	// 	var keys []string
// 	// 	for k := range event.Actor.Attributes {
// 	// 		keys = append(keys, k)
// 	// 	}
// 	// 	sort.Strings(keys)
// 	// 	for _, k := range keys {
// 	// 		v := event.Actor.Attributes[k]
// 	// 		attrs = append(attrs, fmt.Sprintf("%s=%s", k, v))
// 	// 	}
// 	// 	fmt.Fprintf(out, " (%s)", strings.Join(attrs, ", "))
// 	// }
// 	fmt.Fprint(out, "\n")
// 	return nil
// }
//
// func formatEvent(out io.Writer, event types.Request, tmpl *template.Template) error {
// 	defer out.Write([]byte{'\n'})
// 	return tmpl.Execute(out, event)
// }
