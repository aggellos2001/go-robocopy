package gorobocopy

import (
	"io"
	"os/exec"
	"reflect"
	"strconv"

	"github.com/aggellos2001/go-robocopy/flags/aflags"
	"github.com/aggellos2001/go-robocopy/flags/copyflags"
	"github.com/aggellos2001/go-robocopy/flags/dcopyflags"
	"github.com/aggellos2001/go-robocopy/flags/unitflags"
	"github.com/aggellos2001/go-robocopy/types"
)

type Robocopy struct {
	source      string // Specifies the path to the source directory.
	destination string // Specifies the path to the destination directory.
	file        string // Specifies the file or files to be copied. Wildcard characters (* or ?) are supported. If you don't specify this parameter, *.* is used as the default value.

	// Specifies the options to use with the robocopy command, including copy, file, retry, logging, and job options.
	copyOpt       *CopyOptions
	throttlingOpt *CopyFileThrottlingOptions
	fileslOpt     *FileSelectionOptions
	retryOpt      *RetryOptions
	loggingOpt    *LoggingOptions
	jobOpt        *JobOptions
	exitCode      ExitCode
}

type CopyOptions struct {
	// [/s] Copies subdirectories. This option automatically excludes empty directories.
	S bool
	// [/e] Copies subdirectories. This option automatically includes empty directories.
	E bool
	// [/lev:n] Copies only the top n levels of the source directory tree.
	Lev int
	// [/z] Copies files in restartable mode. In restartable mode, should a file copy be interrupted, robocopy can pick up where it left off rather than recopying the entire file.
	Z bool
	// [/b] Copies files in backup mode. In backup mode, robocopy overrides file and folder permission settings (ACLs), which might otherwise block access.
	B bool
	// [/zb] Copies files in restartable mode. If file access is denied, switches to backup mode.
	Zb bool
	// [/j] Copies using unbuffered I/O (recommended for large files).
	J bool
	// [/efsraw] Copies all encrypted files in EFS RAW mode.
	EsfRaw bool
	// [/copy:<copyflags>] Specifies which file properties to copy. The valid values for this option are:
	// D - Data, A - Attributes, T - Time stamps, X - Skip alt data streams, S - NTFS access control list (ACL), O - Owner information, U - Auditing information
	// The default value for the /COPY option is DAT (data, attributes, and time stamps). The X flag is ignored if either /B or /ZB is used.
	Copy copyflags.CopyFlags
	// [/dcopy:<copyflags>] Specifies what to copy in directories. The valid values for this option are:
	// D - Data, A - Attributes, T - Time stamps, E - Extended attribute, X - Skip alt data streams
	// The default value for this option is DA (data and attributes).
	Dcopy dcopyflags.DCopyFlags
	// [/sec] Copies files with security (equivalent to /copy:DATS).
	Sec bool
	// [/copyall] Copies all file information (equivalent to /copy:DATSOU).
	CopyAll bool
	// [/nocopy] Copies no file information (useful with /purge).
	NoCopy bool
	// [/secfix] Fixes file security on all files, even skipped ones.
	SecFix bool
	// [/timfix] Fixes file times on all files, even skipped ones.
	TimFix bool
	// [/purge] Deletes destination files and directories that no longer exist in the source. Using this option with the /e option and a destination directory, allows the destination directory security settings to not be overwritten.
	Purge bool
	// [/mir]  	Mirrors a directory tree (equivalent to /e plus /purge). Using this option with the /e option and a destination directory, overwrites the destination directory security settings.
	Mir bool
	// [/mov] Moves files, and deletes them from the source after they're copied.
	Mov bool
	// [/move] Moves files and directories, and deletes them from the source after they're copied.
	Move bool
	// [/a+:[RASHCNET]]
	// Adds the specified attributes to copied files. The valid values for this option are:
	// R - Read only, A - Archive, S - System, H - Hidden, C - Compressed, N - Not content indexed, E - Encrypted, T - Temporary, O - Offline
	APlus aflags.AFlags
	// [/a-:[RASHCNETO]]
	// Removes the specified attributes from copied files. The valid values for this option are:
	// Everything is the same as [/a+:] with the only addintional value O - Offline
	AMinus aflags.AFlags
	// [/create] Creates a directory tree and zero-length files only.
	Create bool
	// [/fat] Creates destination files by using 8.3 character-length FAT file names only.
	Fat bool
	// [/256] Turns off support for paths longer than 256 characters.
	NoMoreThan256 bool
	// [/mon:n] Monitors the source and runs again when more than n changes are detected.
	Mon int
	// [/mot:m] Monitors the source and runs again in m minutes if changes are detected.
	Mot int
	// [/rh:hhmm-hhmm] Specifies run times when new copies can be started.
	Rh string
	// [/pf] Checks run times on a per file (not per-pass) basis.
	Pf bool
	// [/ipg:n] Specifies the inter-packet gap to free bandwidth on slow lines.
	Ipg int
	// [/sj] Copies junctions (soft-links) to the destination path instead of link targets.
	Sj bool
	// [/sl] Don't follow symbolic links and instead create a copy of the link.
	Sl bool
	// [/mt:n] Creates multi-threaded copies with n threads. n must be an integer between 1 and 128. The default value for n is 8. For better performance, redirect your output using /log option.
	//The /mt parameter can't be used with the /ipg and /efsraw parameters.
	Mt int
	// [/nodcopy] Copies no directory info (the default /dcopy:DA is done).
	Nodcopy bool
	// [/nooffload] Copies files without using the Windows Copy Offload mechanism.
	Nooffload bool
	// [/compress] Requests network compression during file transfer, if applicable.
	Compress bool
	// [/sparse] Enables retaining the sparse state of files during copy.
	Sparse bool
}

func (c *CopyOptions) getCommandArgs() (result []string) {
	if c.S {
		result = append(result, "/s")
	}
	if c.E {
		result = append(result, "/e")
	}
	if c.Lev != 0 {
		result = append(result, "/lev:"+strconv.Itoa(c.Lev))
	}
	if c.Z {
		result = append(result, "/z")
	}
	if c.B {
		result = append(result, "/b")
	}
	if c.Zb {
		result = append(result, "/zb")
	}
	if c.J {
		result = append(result, "/j")
	}
	if c.EsfRaw {
		result = append(result, "/efsraw")
	}
	if c.Copy != 0 {
		result = append(result, "/copy:"+c.Copy.String())
	}
	if c.Dcopy != 0 {
		result = append(result, "/dcopy:"+c.Dcopy.String())
	}
	if c.Sec {
		result = append(result, "/sec")
	}
	if c.CopyAll {
		result = append(result, "/copyall")
	}
	if c.NoCopy {
		result = append(result, "/nocopy")
	}
	if c.SecFix {
		result = append(result, "/secfix")
	}
	if c.TimFix {
		result = append(result, "/timfix")
	}
	if c.Purge {
		result = append(result, "/purge")
	}
	if c.Mir {
		result = append(result, "/mir")
	}
	if c.Mov {
		result = append(result, "/mov")
	}
	if c.Move {
		result = append(result, "/move")
	}
	if c.APlus != 0 {
		result = append(result, "/a+:"+c.APlus.String())
	}
	if c.AMinus != 0 {
		result = append(result, "/a-:"+c.AMinus.String())
	}
	if c.Create {
		result = append(result, "/create")
	}
	if c.Fat {
		result = append(result, "/fat")
	}
	if c.NoMoreThan256 {
		result = append(result, "/256")
	}
	if c.Mon != 0 {
		result = append(result, "/mon:"+strconv.Itoa(c.Mon))
	}
	if c.Mot != 0 {
		result = append(result, "/mot:"+strconv.Itoa(c.Mot))
	}
	if c.Rh != "" {
		result = append(result, "/rh:"+c.Rh)
	}
	if c.Pf {
		result = append(result, "/pf")
	}
	if c.Ipg != 0 {
		result = append(result, "/ipg:"+strconv.Itoa(c.Ipg))
	}
	if c.Sj {
		result = append(result, "/sj")
	}
	if c.Sl {
		result = append(result, "/sl")
	}
	if c.Mt != 0 {
		result = append(result, "/mt:"+strconv.Itoa(c.Mt))
	}
	if c.Nodcopy {
		result = append(result, "/nodcopy")
	}
	if c.Nooffload {
		result = append(result, "/nooffload")
	}
	if c.Compress {
		result = append(result, "/compress")
	}
	if c.Sparse {
		result = append(result, "/sparse")
	}
	return result
}

// These throttling options are used to specify the maximum I/O bandwidth that Robocopy allows to be used in bytes per second. If not specifying in bytes per second, whole numbers can be used if k, m, or g are specified. The minimum I/O bandwidth that is throttled is 524288 bytes even if a lesser value is specified.
type CopyFileThrottlingOptions struct {
	// [/iomaxsize:n[kmg]] The requested max i/o size per read/write cycle in n kilobytes, megabytes, or gigabytes.
	Iomaxsize types.Pair[int, unitflags.UnitFlags]
	// [/iorate:<n>[kmg]] The requested i/o rate in n kilobytes megabytes, or gigabytes per second.
	Iorate types.Pair[int, unitflags.UnitFlags]
	// [/threshold:<n>[kmg]] The file size threshold for throttling in n kilobytes, megabytes, or gigabytes.
	Threshold types.Pair[int, unitflags.UnitFlags]
}

func (cfto *CopyFileThrottlingOptions) getCommandArgs() (result []string) {
	if !reflect.ValueOf(cfto.Iomaxsize).IsZero() {
		result = append(result, "/iomaxsize:"+strconv.Itoa(cfto.Iomaxsize.First)+cfto.Iomaxsize.Second.String())
	}
	if !reflect.ValueOf(cfto.Iorate).IsZero() {
		result = append(result, "/iorate:"+strconv.Itoa(cfto.Iorate.First)+cfto.Iorate.Second.String())
	}
	if !reflect.ValueOf(cfto.Threshold).IsZero() {
		result = append(result, "/threshold:"+strconv.Itoa(cfto.Threshold.First)+cfto.Threshold.Second.String())
	}
	return result
}

type FileSelectionOptions struct {
	// [/a] Copies only files for which the Archive attribute is set.
	A bool
	// [/m] Copies only files for which the Archive attribute is set, and resets the Archive attribute.
	M bool
	// [/ia:[RASHCNETO]] Includes only files for which any of the specified attributes are set. The valid values for this option are:
	// R - Read only, A - Archive, S - System, H - Hidden, C - Compressed, N - Not content indexed, E - Encrypted, T - Temporary, O - Offline
	Ia aflags.AFlags
	// [//xa:[RASHCNETO]] Excludes files for which any of the specified attributes are set. The valid values for this option are the same as the [/ia] command.
	Xa aflags.AFlags
	// [/xf <filename>[...]] Excludes files that match the specified names or paths. Wildcard characters (* and ?) are supported.
	Xf []string
	// [/xd <directory>[...]] Excludes directories that match the specified names and paths.
	Xd []string
	// [/xc] Excludes existing files with the same timestamp, but different file sizes.
	Xc bool
	// [/xn] Source directory files newer than the destination are excluded from the copy.
	Xn bool
	// [/xo] Source directory files older than the destination are excluded from the copy.
	Xo bool
	// [/xx] Excludes extra files and directories present in the destination but not the source. Excluding extra files won't delete files from the destination.
	Xx bool
	// [/xl] Excludes "lonely" files and directories present in the source but not the destination. Excluding lonely files prevents any new files from being added to the destination.
	Xl bool
	// [/im] Include modified files (differing change times).
	Im bool
	// [/is] Includes the same files. Same files are identical in name, size, times, and all attributes.
	Is bool
	// [/it] Includes "tweaked" files. Tweaked files have the same name, size, and times, but different attributes.
	It bool
	// [/max:n] Specifies the maximum file size (to exclude files bigger than n bytes).
	Max int
	// [/min:n] Specifies the minimum file size (to exclude files smaller than n bytes).
	Min int
	// [/maxage:n] Specifies the maximum file age (to exclude files older than n days or date).
	Maxage int
	// [/minage:n] Specifies the minimum file age (exclude files newer than n days or date).
	Minage int
	// [/maxlad:n] Specifies the maximum last access date (excludes files unused since n).
	Maxlad int
	// [/minlad:n] Specifies the minimum last access date (excludes files used since n) If n is less than 1900, n specifies the number of days. Otherwise, n specifies a date in the format YYYYMMDD.
	Minlad int
	// [/xj] Excludes junction points, which are normally included by default.
	Xj bool
	// [/fft] Assumes FAT file times (two-second precision).
	Fft bool
	// [/dst] Compensates for one-hour DST time differences.
	Dst bool
	// [/xjd] Excludes junction points for directories.
	Xjd bool
	// [/xjf] Excludes junction points for files.
	Xjf bool
}

func (fso *FileSelectionOptions) getCommandArgs() (result []string) {
	if fso.A {
		result = append(result, "/a")
	}
	if fso.M {
		result = append(result, "/m")
	}
	if fso.Ia != 0 {
		result = append(result, "/ia:"+fso.Ia.String())
	}
	if fso.Xa != 0 {
		result = append(result, "/xa:"+fso.Xa.String())
	}
	if len(fso.Xf) != 0 {
		result = append(result, "/xf")
		result = append(result, fso.Xf...)
	}
	if len(fso.Xd) != 0 {
		result = append(result, "/xd")
		result = append(result, fso.Xd...)
	}
	if fso.Xc {
		result = append(result, "/xc")
	}
	if fso.Xn {
		result = append(result, "/xn")
	}
	if fso.Xo {
		result = append(result, "/xo")
	}
	if fso.Xx {
		result = append(result, "/xx")
	}
	if fso.Xl {
		result = append(result, "/xl")
	}
	if fso.Im {
		result = append(result, "/im")
	}
	if fso.Is {
		result = append(result, "/is")
	}
	if fso.It {
		result = append(result, "/it")
	}
	if fso.Max != 0 {
		result = append(result, "/max:"+strconv.Itoa(fso.Max))
	}
	if fso.Min != 0 {
		result = append(result, "/min:"+strconv.Itoa(fso.Min))
	}
	if fso.Maxage != 0 {
		result = append(result, "/maxage:"+strconv.Itoa(fso.Maxage))
	}
	if fso.Minage != 0 {
		result = append(result, "/minage:"+strconv.Itoa(fso.Minage))
	}
	if fso.Maxlad != 0 {
		result = append(result, "/maxlad:"+strconv.Itoa(fso.Maxlad))
	}
	if fso.Minlad != 0 {
		result = append(result, "/minlad:"+strconv.Itoa(fso.Minlad))
	}
	if fso.Xj {
		result = append(result, "/xj")
	}
	if fso.Fft {
		result = append(result, "/fft")
	}
	if fso.Dst {
		result = append(result, "/dst")
	}
	if fso.Xjd {
		result = append(result, "/xjd")
	}
	if fso.Xjf {
		result = append(result, "/xjf")
	}
	return result
}

type RetryOptions struct {
	// [/r:<n>] Specifies the number of retries on failed copies. The default value of n is 1,000,000 (one million retries).
	R int
	// [/w:<n>] Specifies the wait time between retries, in seconds. The default value of n is 30 (wait time 30 seconds).
	W int
	// [/reg] Saves the values specified in the /r and /w options as default settings in the registry.
	Reg bool
	// [/tbd] Specifies that the system waits for share names to be defined (retry error 67).
	Tbd bool
	// [/lfsm] Operate in low free space mode that enables copy, pause, and resume (see Remarks).
	Lfsm bool
	// [/lfsm:<n>[kmg]] Specifies the floor size in n kilobytes, megabytes, or gigabytes.
	LfsmSize types.Pair[int, unitflags.UnitFlags]
}

func (ropt *RetryOptions) getCommandArgs() (result []string) {
	if ropt.R != 0 {
		result = append(result, "/r:"+strconv.Itoa(ropt.R))
	}
	if ropt.W != 0 {
		result = append(result, "/w:"+strconv.Itoa(ropt.W))
	}
	if ropt.Reg {
		result = append(result, "/reg")
	}
	if ropt.Tbd {
		result = append(result, "/tbd")
	}
	if ropt.Lfsm {
		result = append(result, "/lfsm")
	}
	if !reflect.ValueOf(ropt.LfsmSize).IsZero() {
		result = append(result, "/lfsm:"+strconv.Itoa(ropt.LfsmSize.First)+ropt.LfsmSize.Second.String())
	}
	return result
}

type LoggingOptions struct {
	// [/l] Specifies that files are to be listed only (and not copied, deleted, or time stamped).
	L bool
	// [/x] Reports all extra files, not just the ones that are selected.
	X bool
	// [/v] Produces verbose output, and shows all skipped files.
	V bool
	// [/ts] Includes source file time stamps in the output.
	Ts bool
	// [/fp] Includes the full path names of the files in the output.
	Fp bool
	// [/bytes] Prints sizes, as bytes.
	Bytes bool
	// [/ns] Specifies that file sizes aren't to be logged.
	Ns bool
	// [/nc] Specifies that file classes aren't to be logged.
	Nc bool
	// [/nfl] Specifies that file names aren't to be logged.
	Nfl bool
	// [/ndl] Specifies that directory names aren't to be logged.
	Ndl bool
	// [/np] Specifies that the progress of the copying operation (the number of files or directories copied so far) won't be displayed.
	Np bool
	// [/eta] Shows the estimated time of arrival (ETA) of the copied files.
	Eta bool
	// [/log:logfile] Writes the status output to the log file (overwrites the existing log file).
	Log string
	// [/log+:logfile] Writes the status output to the log file (appends the output to the existing log file).
	LogPlus string
	// [/unilog:logfile] Writes the status output to the log file as unicode text (overwrites the existing log file).
	UniLog string
	// [/unilog+:logfile] Writes the status output to the log file as Unicode text (appends the output to the existing log file).
	UniLogPlus string
	// [/tee] Writes the status output to the console window, and to the log file.
	Tee bool
	// [/njh] Specifies that there's no job header.
	Njh bool
	// [/njs] Specifies that there's no job summary.
	Njs bool
	// [/unicode] Displays the status output as unicode text.
	Unicode bool
}

func (lopt *LoggingOptions) getCommandArgs() (result []string) {
	if lopt.L {
		result = append(result, "/l")
	}
	if lopt.X {
		result = append(result, "/x")
	}
	if lopt.V {
		result = append(result, "/v")
	}
	if lopt.Ts {
		result = append(result, "/ts")
	}
	if lopt.Fp {
		result = append(result, "/fp")
	}
	if lopt.Bytes {
		result = append(result, "/bytes")
	}
	if lopt.Ns {
		result = append(result, "/ns")
	}
	if lopt.Nc {
		result = append(result, "/nc")
	}
	if lopt.Nfl {
		result = append(result, "/nfl")
	}
	if lopt.Ndl {
		result = append(result, "/ndl")
	}
	if lopt.Np {
		result = append(result, "/np")
	}
	if lopt.Eta {
		result = append(result, "/eta")
	}
	if lopt.Log != "" {
		result = append(result, "/log:"+lopt.Log)
	}
	if lopt.LogPlus != "" {
		result = append(result, "/log+:"+lopt.LogPlus)
	}
	if lopt.UniLog != "" {
		result = append(result, "/unilog:"+lopt.UniLog)
	}
	if lopt.UniLogPlus != "" {
		result = append(result, "/unilog+:"+lopt.UniLogPlus)
	}
	if lopt.Tee {
		result = append(result, "/tee")
	}
	if lopt.Njh {
		result = append(result, "/njh")
	}
	if lopt.Njs {
		result = append(result, "/njs")
	}
	if lopt.Unicode {
		result = append(result, "/unicode")
	}
	return result
}

type JobOptions struct {
	// [/job:jobname] Specifies that parameters are to be derived from the named job file. To run /job:jobname, you must first run the /save:jobname parameter to create the job file.
	Job string
	// [/save:jobname] Specifies that parameters are to be saved to the named job file. This must be ran before running /job:jobname. All copy, retry, and logging options must be specified before this parameter.
	Save string
	// [/quit] Quits after processing command line (to view parameters).
	Quit bool
	// [/nosd] Indicates that no source directory is specified.
	Nosd bool
	// [/nodd] Indicates that no destination directory is specified.
	Nodd bool
	// [/if] Includes the specified files.
	If bool
}

func (jopt *JobOptions) getCommandArgs() (result []string) {
	if jopt.Job != "" {
		result = append(result, "/job:"+jopt.Job)
	}
	if jopt.Quit {
		result = append(result, "/quit")
	}
	if jopt.Nosd {
		result = append(result, "/nosd")
	}
	if jopt.Nodd {
		result = append(result, "/nodd")
	}
	if jopt.If {
		result = append(result, "/if")
	}
	if jopt.Save != "" {
		result = append(result, "/save:"+jopt.Save)
	}
	return result
}

// NewRobocopy returns a new robocopy instance with the default options applied.
func NewRobocopy(sourceDir, destinationDir, file string) *Robocopy {
	return &Robocopy{
		source:      sourceDir,
		destination: destinationDir,
		file:        file,
	}
}

func (r *Robocopy) SetCopyOptions(opts *CopyOptions) {
	r.copyOpt = opts
}

func (r *Robocopy) SetThrottlingOptions(opts *CopyFileThrottlingOptions) {
	r.throttlingOpt = opts
}

func (r *Robocopy) SetFileSelectionOptions(opts *FileSelectionOptions) {
	r.fileslOpt = opts
}

func (r *Robocopy) SetRetryOptions(opts *RetryOptions) {
	r.retryOpt = opts
}

func (r *Robocopy) SetLoggingOptions(opts *LoggingOptions) {
	r.loggingOpt = opts
}

func (r *Robocopy) SetJobOptions(opts *JobOptions) {
	r.jobOpt = opts
}

// Returns the command arguments for the robocopy command in
// the form of a string slice that can be used with exec.Command.
func (r *Robocopy) GetCommandArgs() (command []string) {
	command = append(command, r.source)
	command = append(command, r.destination)
	command = append(command, r.file)
	if r.copyOpt != nil {
		command = append(command, r.copyOpt.getCommandArgs()...)
	}
	if r.throttlingOpt != nil {
		command = append(command, r.throttlingOpt.getCommandArgs()...)
	}
	if r.fileslOpt != nil {
		command = append(command, r.fileslOpt.getCommandArgs()...)
	}
	if r.retryOpt != nil {
		command = append(command, r.retryOpt.getCommandArgs()...)
	}
	if r.loggingOpt != nil {
		command = append(command, r.loggingOpt.getCommandArgs()...)
	}
	if r.jobOpt != nil {
		command = append(command, r.jobOpt.getCommandArgs()...)
	}
	return command
}

// Returns a ready-to-use exec.Cmd instance for the robocopy command.
// It is up to the caller to run the command and handle the output.
func (r *Robocopy) GetCommand() *exec.Cmd {
	return exec.Command("robocopy", r.GetCommandArgs()...)
}

// Handles running the command and populating the exit code.
// You can set the stdin, stdout, and stderr of the command.
// If you set to nil, they will be set to the nul device (os.DevNull).
// You can check the exit code using the GetExitCode() function.
func (r *Robocopy) Run(stdin io.Reader, stdout, stderr io.Writer) {
	cmd := r.GetCommand()
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	cmd.Run()
	r.exitCode = ExitCode(cmd.ProcessState.ExitCode())
}

type ExitCode int

const (
	// No files were copied. No failure was encountered. No files were mismatched.
	// The files already exist in the destination directory; therefore, the copy operation was skipped.
	AlreadyExist ExitCode = iota
	// All files were copied successfully.
	AllFilesCopied
	// There are some additional files in the destination directory that aren't present in the source directory. No files were copied.
	AdditionalFilesOnDest
	// Some files were copied. Additional files were present. No failure was encountered.
	SomeFilesCopied
	// For some reason exit code 4 is skipped according to Microsoft's documentation.
	_
	// Some files were copied. Some files were mismatched. No failure was encountered.
	SomeFilesMismatched
	// Additional files and mismatched files exist. No files were copied and no failures were encountered meaning that the files already exist in the destination directory.
	AdditionalAndMismatchedFiles
	// Files were copied, a file mismatch was present, and additional files were present.
	FilesCopiedMismatchedAndAdditional
	// Several files didn't copy.
	SeveralFilesDidntCopy
)

func (r *Robocopy) GetExitCode() ExitCode {
	return r.exitCode
}

func (e ExitCode) String() string {
	switch e {
	case AlreadyExist:
		return "No files were copied. No failure was encountered. No files were mismatched. The files already exist in the destination directory; therefore, the copy operation was skipped."
	case AllFilesCopied:
		return "All files were copied successfully."
	case AdditionalFilesOnDest:
		return "There are some additional files in the destination directory that aren't present in the source directory. No files were copied."
	case SomeFilesCopied:
		return "Some files were copied. Additional files were present. No failure was encountered."
	case SomeFilesMismatched:
		return "Some files were copied. Some files were mismatched. No failure was encountered."
	case AdditionalAndMismatchedFiles:
		return "Additional files and mismatched files exist. No files were copied and no failures were encountered meaning that the files already exist in the destination directory."
	case FilesCopiedMismatchedAndAdditional:
		return "Files were copied, a file mismatch was present, and additional files were present."
	case SeveralFilesDidntCopy:
		return "Several files didn't copy."
	}
	return "Unknown exit code"
}
