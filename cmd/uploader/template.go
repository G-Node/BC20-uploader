package main

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
		<title>BC20 Poster Submission</title>
		<meta name="twitter:card" content="summary" />
		<meta name="twitter:site" content="@nncn_germany" />
		<meta name="twitter:title" content="BC20 Poster Submission"/>
		<meta name="twitter:description" content="BC20 Poster Submission"/>
		<meta name="twitter:image" content="/assets/favicon.png" />
	</head>
	<body>
		<div class="full height">
			<div class="following bar light">
				<div class="ui container">
					<div class="ui grid">
						<div class="column">
							<div class="ui top secondary menu">
								<a class="item brand" href="https://bc20-posters.g-node.org">
									<img class="ui mini image" src="/assets/favicon.png">
									<a class="item" href="http://www.bernstein-conference.de/">Conference Website</a>
									<a class="item" href="mailto:bernstein.conference@fz-juelich.de">Contact</a>
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
					<a href="http://www.g-node.org"><img class="ui mini footericon" src="https://projects.g-node.org/assets/gnode-bootstrap-theme/1.2.0-snapshot/img/gnode-icon-50x50-transparent.png"/>© G-Node, 2016-2020</a>
					<a href="https://gin.g-node.org/G-Node/Info/wiki/about">About</a>
					<a href="https://gin.g-node.org/G-Node/Info/wiki/imprint">Imprint</a>
					<a href="https://gin.g-node.org/G-Node/Info/wiki/contact">Contact</a>
					<a href="https://gin.g-node.org/G-Node/Info/wiki/Terms+of+Use">Terms of Use</a>
					<a href="https://gin.g-node.org/G-Node/Info/wiki/Datenschutz">Datenschutz</a>
				</div>
				<div class="ui center links item brand footertext">
					<span>Powered by:      <a href="https://github.com/gogits/gogs"><img class="ui mini footericon" src="https://gin.g-node.org/img/gogs.svg"/></a>         </span>
					<span>Hosted by:       <a href="http://neuro.bio.lmu.de"><img class="ui mini footericon" src="https://gin.g-node.org/img/lmu.png"/></a>          </span>
					<span>Funded by:       <a href="http://www.bmbf.de"><img class="ui mini footericon" src="https://gin.g-node.org/img/bmbf.png"/></a>         </span>
					<span>Registered with: <a href="http://doi.org/10.17616/R3SX9N"><img class="ui mini footericon" src="https://gin.g-node.org/img/re3.png"/></a>          </span>
					<span>Recommended by:  <a href="https://www.nature.com/sdata/policies/repositories#neurosci"><img class="ui mini footericon" src="https://gin.g-node.org/img/sdatarecbadge.jpg"/><a href="https://journals.plos.org/plosone/s/data-availability#loc-neuroscience"><img class="ui mini footericon" src="https://gin.g-node.org/img/sm_plos-logo-sm.png"/></a></span>
				</div>
			</div>
		</footer>
	</body>
</html>
{{ end }}
`

const Form = `
{{ define "content" }}
			<div class="ginform">
				<div class="ui middle very relaxed page grid">
					<div class="column">
						<form class="ui form" enctype="multipart/form-data" action="/submit" method="post">
							<input type="hidden" name="_csrf" value="">
							<h3 class="ui top attached header">
								BC20 Poster Submission
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
								<span class="help">Link to short poster presentation video</span>
							</div>
							<div class="inline required field ">
								<label for="passcode">Passcode</label>
								<input type="password" id="passcode" name="passcode" value="" autofocus required>
								<span class="help">You should have received a passcode in the instruction email</span>
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
const SuccessTmpl = `
{{ define "content" }}
			<div class="home middle very relaxed page grid" id="main">
				<div class="ui container wide centered column doi">
					<div class="column center">
						<h1>BC20 poster upload service</h1>
					</div>
					<div class="ui info message" id="infotable">
						<div id="infobox">
							<p>The following <strong>preview</strong> shows the information that will appear in the poster gallery alongside your poster.</p>
							<p>Please review it carefully and <strong><a href="mailto:bernstein.conference@fz-juelich.de">contact us</a></strong> if there are any issues.</p>
						</div>
					</div>
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

const FailureTmpl = `
{{ define "content" }}
			<div class="home middle very relaxed page grid" id="main">
				<div class="ui container wide centered column doi">
					<div class="column center">
						<h1>BC20 poster upload service</h1>
					</div>
					<div class="ui error message" id="infotable">
						<div id="infobox">
							<p>The poster submission failed.<p>

							<p>{{.Message}}</p>

							<p>Please <strong><a href="mailto:bernstein.conference@fz-juelich.de">contact us</a></strong> if there are any issues. <a href="/">Click here</a> to return to the form and try again.</p>
						</div>
					</div>
					<hr>
				</div>
			</div>
		</div>
{{end}}
`

// vim: ft=gohtmltmpl
