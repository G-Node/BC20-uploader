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
		<title>Project creator</title>
		<meta name="twitter:card" content="summary" />
		<meta name="twitter:site" content="@gnode" />
		<meta name="twitter:title" content="GIN Valid"/>
		<meta name="twitter:description" content="G-Node GIN Validation service"/>
		<meta name="twitter:image" content="/assets/favicon.png" />
	</head>
	<body>
		<div class="full height">
			<div class="following bar light">
				<div class="ui container">
					<div class="ui grid">
						<div class="column">
							<div class="ui top secondary menu">
								<a class="item brand" href="https://gin.g-node.org/">
									<img class="ui mini image" src="/assets/favicon.png">
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
					<a href="http://www.g-node.org"><img class="ui mini footericon" src="https://projects.g-node.org/assets/gnode-bootstrap-theme/1.2.0-snapshot/img/gnode-icon-50x50-transparent.png"/>© G-Node, 2016-2019</a>
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

// Form and Job view page template
const Form = `
{{ define "content" }}
			<div class="ginform">
				<div class="ui middle very relaxed page grid">
					<div class="column">
						<form class="ui form" enctype="multipart/form-data" action="/submit" method="post">
							<input type="hidden" name="_csrf" value="">
							<h3 class="ui top attached header">
								(DEMO) BC20 poster session upload
							</h3>
							<div class="ui attached segment">
								<div class="inline required field ">
									<label for="poster">Poster (PDF)</label>
									<input type="file" id="poster" name="poster" accept="application/pdf" required>
									<span class="help">Poster or slides</span>
								</div>
								{{if .videos}}
									<div class="inline field ">
										<label for="video">Video</label>
										<input type="file" id="video" name="video" accept="video/*">
										<span class="help">Short poster presentation video</span>
									</div>
								{{end}}
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

// vim: ft=gohtmltmpl
