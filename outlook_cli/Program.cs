using System;
using Microsoft.Office.Interop.Outlook;
using CommandLine;

namespace outlook_cli
{

    class Program
    {
        public class Options
        {
            [Value(0, MetaName = "Command", Required = true, HelpText = "Command to execute")]
            public string Command { get; set; }

            [Option('s', "subject", Required = false, HelpText = "Subject of event")]
            public string Subject { get; set; }

            [Option('b', "body", Required = false, HelpText = "Body of event")]
            public string Body { get; set; }

            [Option('l', "location", Required = false, HelpText = "Location of event")]
            public string Location { get; set; }

            [Option('d', "startDate", Required = false, HelpText = "Start date of event")]
            public string StartDate { get; set; }

            [Option('t', "startTime", Required = false, HelpText = "Start time of event")]
            public string StartTime { get; set; }

            [Option('e', "endTime", Required = false, HelpText = "End time of event")]
            public string EndTime { get; set; }

            [Option('r', "dryRun", Required = false, HelpText = "Don't do any changes, just print")]
            public bool DryRun { get; set; }

            [Option('f', "filter", Required = false, HelpText = "Filter for list command")]
            public string Filter { get; set; }
        }

        static void Main(string[] args)
        {
            Parser.Default.ParseArguments<Options>(args).WithParsed(run);
        }

        static void run(Options opts)
        {
            switch (opts.Command)
            {
                case "add":
                    addEvent(opts);
                    break;
                case "remove":
                    removeEvent(opts);
                    break;
                case "list":
                    listEvents(opts);
                    break;
                default:
                    Console.WriteLine("Unknown command " + opts.Command);
                    Environment.Exit(1);
                    break;
            }
        }

        static void addEvent(Options opts)
        {
            assertRequired("Subject", opts.Subject);
            assertRequired("Start date", opts.StartDate);

            Application outlookApp = new Application();
            AppointmentItem oAppointment = (AppointmentItem)outlookApp.CreateItem(OlItemType.olAppointmentItem);

            oAppointment.Subject = opts.Subject;
            oAppointment.Body = opts.Body;

            oAppointment.Start = getFullTime(opts.StartDate, opts.StartTime);
            if (!string.IsNullOrEmpty(opts.EndTime))
            {
                oAppointment.End = getFullTime(opts.StartDate, opts.EndTime);
            }
            else
            {
                oAppointment.AllDayEvent = true;
            }

            
            oAppointment.ReminderSet = true;
            oAppointment.ReminderMinutesBeforeStart = 15;
            oAppointment.Sensitivity = OlSensitivity.olPrivate;
            oAppointment.BusyStatus = OlBusyStatus.olFree;
            oAppointment.Location = opts.Location;

            if (opts.DryRun)
            {
                Console.WriteLine("Would create event:");
                Console.WriteLine(apptToString(oAppointment));
            }
            else
            {
                oAppointment.Save();
            }
        }

        static void removeEvent(Options opts)
        {
            assertRequired("Subject", opts.Subject);

            Application outlookApp = new Application();
            Folder calFolder = outlookApp.Session.GetDefaultFolder(OlDefaultFolders.olFolderCalendar) as Folder;
            var appt = calFolder.Items.Find("[Subject] = '" + opts.Subject  +"'") as AppointmentItem;
            if (appt == null)
            {
                Console.WriteLine("Could not find event with subject " + opts.Subject);
                Environment.Exit(1);
            }

            if (opts.DryRun)
            {
                Console.WriteLine("Would delete event:");
                Console.WriteLine(apptToString(appt));
            }
            else
            {
                appt.Delete();
            }
        }

        static void listEvents(Options opts)
        {
            assertRequired("Filter", opts.Filter);

            Application outlookApp = new Application();
            Folder calFolder = outlookApp.Session.GetDefaultFolder(OlDefaultFolders.olFolderCalendar) as Folder;

            var found = calFolder.Items.Restrict(opts.Filter);
            foreach (AppointmentItem appt in found)
            {
                Console.WriteLine(apptToCsv(appt));
            }
        }

        static DateTime getFullTime(string dateString, string timeString)
        {
            var start = DateTime.Parse(dateString);
            if (!string.IsNullOrEmpty(timeString))
            {
                var time = TimeSpan.Parse(timeString);
                start = start.Date + time;
            }
            return start;
        }

        static string apptToString(AppointmentItem oAppointment)
        {
            return "Subject: " + oAppointment.Subject
                + "\nStart: " + oAppointment.Start
                + "\nEnd: " + oAppointment.End
                + "\nAllDayEvent: " + oAppointment.AllDayEvent;
        }

        static string apptToCsv(AppointmentItem oAppointment)
        {
            return oAppointment.Start +
                ";" + oAppointment.End +
                ";" + oAppointment.AllDayEvent +
                ";" + oAppointment.Subject;
        }

        static void assertRequired(string name, string val)
        {
            if (string.IsNullOrEmpty(val))
            {
                Console.WriteLine(name + " required");
                Environment.Exit(1);
            }
        }
    }
}
