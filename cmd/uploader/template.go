package main

// Layout is the base Uploader template providing header and footer.
// If required adjust the main poster gallery URL and the copyright
// date before deploying for a new conference.
const Layout = `
{{ define "layout" }}
<html>
	<!DOCTYPE html>
	<head data-suburl="">
		<link rel="shortcut icon" href="/assets/favicon.png" />
		<link rel="stylesheet" href="/assets/font-awesome-4.6.3/css/font-awesome.min.css">
		<link rel="stylesheet" href="/assets/octicons-4.3.0/octicons.min.css">
		<link rel="stylesheet" href="/assets/semantic-2.3.1.min.css">
		<link rel="stylesheet" href="/assets/gogs.css">
		<link rel="stylesheet" href="/assets/custom.css">
		<title>Bernstein Conference Poster Submission</title>
		<meta name="twitter:card" content="summary" />
		<meta name="twitter:site" content="@nncn_germany" />
		<meta name="twitter:title" content="Bernstein Conference Poster Submission"/>
		<meta name="twitter:description" content="Bernstein Conference Poster Submission"/>
		<meta name="twitter:image" content="/assets/favicon.png" />
	</head>
	<body>
		<div class="full height">
			<div class="following bar light">
				<div class="ui container">
					<div class="ui grid">
						<div class="column">
							<div class="ui top secondary menu">
								<a class="item brand" href="https://posters.bc.g-node.org">
									<img class="ui mini image" src="/assets/favicon.png">
									<a class="item" href="{{ .conferencepageurl }}">Conference Website</a>
									<a class="item" href="mailto:{{ .supportemail }}">Contact</a>
								</a>
							</div>
						</div>
					</div>
				</div>
			</div>
			{{ template "content" . }}
		</div>
		<footer>
			<div class="ui container">
				<div class="ui center links item brand footertext">
					<a href="http://www.g-node.org">
						<img class="ui mini footericon" 
							 src="https://projects.g-node.org/assets/gnode-bootstrap-theme/1.2.0-snapshot/img/gnode-icon-50x50-transparent.png"/>
						© 2020-2022 G-Node
					</a>
					<a href="https://bc.g-node.org/G-Node/Info/wiki/Terms+of+Use">Terms of Use</a>
					<a href="https://bc.g-node.org/G-Node/Info/wiki/Datenschutz">Datenschutz</a>
					<a href="https://bc.g-node.org/G-Node/Info/wiki/imprint">Imprint</a>
				</div>
			</div>
		</footer>
	</body>
</html>
{{ end }}
`

// Form is the poster and video/video link upload submission form template.
// The video file upload field is unavailable if the "Videos" item
// in the servers config is set to 'false'.
// The upload form is replaced by a "submission closed" page, if the
// provided 'submission' variable is 'false'.
const Form = `
{{ define "content" }}
	{{ if .submission }}
			<!-- Poster and video link submission form -->
			<div class="body">
				<div class="ui middle very relaxed page grid">
					<div class="column">
					</div>
				</div>
			</div>
			<div class="ginform">
				<div class="ui middle very relaxed page grid">
					<div class="column">
						<form class="ui form" enctype="multipart/form-data" action="/submit" method="post">
							<input type="hidden" name="_csrf" value="">
							<p>Please upload your PDF and video URL by <strong>{{ .closedtext }}</strong> using the form below.
							You have received an upload key in the instruction email.</p>
							<p>You can access the form and re-upload your poster and URL until the deadline.</p>
							<p><strong>Please note: posters sent via email will not be considered.</strong></p>
							<br/>
							<!-- previous video upload version; might come back in the future -->
							<!--
							<p>If you prefer to have your pre-recorded video hosted by us on the Bernstein Conference Vimeo channel, 
								rather than an individual solution, we offer the following alternative:
							<ul>
								<li>Please upload your video in MP4-format here, by {{ .closedtextvid }}: 
								<a href="{{ .viduploadurl }}">{{ .viduploadurl }}</a></li>
								<li>Label your video-file: <code>yourposter#_lastname_video</code></li>
							</ul>
							</p>
							-->
							<p>If you wish to have the pre-recorded video hosted on the <a href="https://vimeo.com/bernsteinnetwork">
							Bernstein Conference Vimeo channel</a> you must upload your video 
							<b>before {{ .closedtextvid }}</b>! 
							Upload your video file here: <a href="{{ .viduploadurl }}">{{ .viduploadurl }}</a>.
							<br>
							Preferred format is *.mp4. File names <b>must</b> follow the naming scheme <code>AbstractNumber_FirstAuthor</code>. 
							<br>
							<b>Other file names will not be considered</b> and authors cannot be informed.
							</p>
							<h3 class="ui top attached header">
								Bernstein Conference Poster Submission Form
							</h3>
							<div class="ui attached segment">
								<div class="inline required field">
									<label for="poster">Poster (PDF)</label>
									<input type="file" id="poster" name="poster" accept="application/pdf" required>
									<span class="help">Poster or slides</span>
								</div>
								{{if .videos}}
									<div class="inline field">
										<label for="video">Video</label>
										<input type="file" id="video" name="video" accept="video/*">
										<span class="help">Short poster presentation video</span>
									</div>
								{{end}}
								<div class="inline field">
									<label for="video_url">Video URL</label>
									<input type="url" id="video_url" name="video_url">
									<span class="help">Link to short self-hosted presentation</span>
								</div>
								<div class="inline required field ">
									<label for="passcode">Upload key</label>
									<input type="password" id="passcode" name="passcode" value="" autofocus required>
									<span class="help">You have received an upload key in the instruction email</span>
								</div>
								<div class="inline field">
									<label></label>
									<button class="ui green button">Submit</button>
								</div>
							</div>
						</form>
					</div>
				</div>
			</div>
	{{ else }}
			<!-- Closed poster and video link submission-->
			<div class="ui container">
				<div class="jumbotron">
					<div class="page-header">
						<h1>Bernstein Conference Poster Submission</h1>
					</div>

					<a href="{{ .conferencepageurl }}">
						<img class="conference-banner img-responsive img-rounded" src="/assets/BC_online_header.jpeg" alt="Conference Logo">
					</a>

					<br>
					<div class="jumbo-small center">
						<p>Poster and video submission is <b class="red">closed</b>.<br></p>
					</div>
					<br>

					<div class="jumbo-small">
						<p>Each year the Bernstein Network invites the international computational neuroscience community to the annual 
							Bernstein Conference for intensive scientific exchange. It has established itself as one of the most renown 
							conferences worldwide in this field, attracting students, postdocs and PIs from around the world to meet and 
							discuss new scientific discoveries.<br></p>
					</div>
				</div>
			</div>
	{{ end }}
{{ end }}
`

// SuccessTmpl is the page displayed after a successful poster content upload.
// It displays an overview of the posters metadata and links to the uploaded
// content.
const SuccessTmpl = `
{{ define "content" }}
			<div class="home middle very relaxed page grid" id="main">
				<div class="ui container wide centered column doi">
					<div class="column center">
						<h1>Bernstein Conference Poster Submission Success</h1>
					</div>

					<div class="ui info message" id="infotable">
						<div id="infobox">
							<p>Your upload was <strong>successful!</strong></p>
							<p>The following <strong>preview</strong> shows the information that will appear in the poster gallery alongside your poster.</p>
							<p>Please review it carefully and <strong><a href="mailto:{{.supportemail}}">contact us</a></strong> 
								if there are any issues.</p>
						</div>
					</div>
					<div><b>NOTE: Please print this page or save the following for verification. 
					You may be asked to produce the following key to verify your upload.</b></div>
					<div>Poster upload verification: <code>{{.PosterHash}}</code></div>
					<hr>
					{{with .UserData}}
					<div class="doi title">
						<h1>{{.Title}}</h1>
						{{.Authors}}
						<p><strong>Session {{.Session}}</strong> | {{.AbstractNumber}} | {{.Topic}}</p>
					</div>
					<hr>

					<h3>Abstract</h3>
					<p>{{.Abstract}}</p>
					{{end}}

					<div><a href="{{.PDFPath}}">Poster PDF</a> (click to review)</div>
					{{if .VideoURL}}
						<div><a href="{{.VideoURL}}">{{.VideoURL}}</a>: Poster presentation video</div>
					{{end}}
					<hr>
				</div>
			</div>
		</div>
{{end}}
`

// FailureTmpl is the page displayed if a poster upload or a whitelist email
// upload has failed.
const FailureTmpl = `
{{ define "content" }}
			<div class="home middle very relaxed page grid" id="main">
				<div class="ui container wide centered column doi">
					<div class="column center">
						<h1>Bernstein Conference Poster Submission</h1>
					</div>
					<div class="ui error message" id="infotable">
						<div id="infobox">
							<p>The submission failed.<p>

							<p>{{.Message}}</p>

							<p>Please <strong><a href="mailto:{{ .supportemail }}}">contact us</a></strong> 
							if there are any issues. <a href="/">Click here</a> to return to the form and try again.</p>
						</div>
					</div>
					<hr>
				</div>
			</div>
		</div>
{{end}}
`

// EmailFormTmpl is the form displayed to add email addresses to
// the whitelist file.
const EmailFormTmpl = `
{{ define "content" }}
<div class="body">
	<div class="ui middle very relaxed page grid">
		<div class="column"></div>
	</div>
</div>
<div class="ginform">
	<div class="ui middle very relaxed page grid">
		<div class="column">
			<form class="ui form" method='post' action='/submitemail'>
				<h3 class="ui top attached header">
					Bernstein Conference whitelist email address upload form
				</h3>
				<div class="ui attached segment">
					<div class="inline required field">
						<label for='content'>Email addresses</label>
						<textarea required name='content' id='content'></textarea>
						<span class="help">Email addresses can be separated by comma, semicolon, space, tab or newline.
						You can always upload a full list, only new addresses are added.</span>
					</div>
					<div class="inline required field">
						<label for='password'>Password</label>
						<input required type='password' name='password' id='password'>
					</div>
					<div class="inline field">
						<label></label>
						<button class="ui green button">Submit</button>
					</div>
				</div>
			</form>
		</div>
	</div>
</div>
{{ end }}
`

// EmailSubmitTmpl is the white list email upload success page.
const EmailSubmitTmpl = `
{{ define "content" }}
<div class="ui container">
	<p></p>
	<h1>Upload received</h1>
	<div class="ui dividing header"></div>
	<p>Whitelist email addresses have been uploaded.</p>
	<p><a href='/uploademail'>Back to the email upload form</a></p>
</div>
{{ end }}
`

// EmailFailTmpl is the page displayed when an error other than 'invalid password'
// occurred during the email whitelist upload.
const EmailFailTmpl = `
{{ define "content" }}
<div class="home middle very relaxed page grid" id="main">
	<div class="ui container wide centered column doi">
		<div class="column center">
			<h1>Bernstein Conference whitelist email upload</h1>
		</div>
		<div class="ui error message" id="infotable">
			<div id="infobox">
				<p>The upload has failed.<p>

				<p>{{.Message}}</p>

				<p><a href="/uploademail">Click here</a> to return to the upload form and try again.</p>
			</div>
		</div>
		<hr>
	</div>
</div>
{{ end }}
`

// vim: ft=gohtmltmpl
