<!DOCTYPE html>
<html>
<head>
  <title>{{.Title}} - 李强的博客</title>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <link href="https://cdn.bootcss.com/bootstrap/3.3.7/css/bootstrap.min.css" rel="stylesheet">
  <link href="https://cdn.bootcss.com/github-markdown-css/2.10.0/github-markdown.min.css" rel="stylesheet">
  <link href="https://cdn.bootcss.com/highlight.js/9.12.0/styles/github.min.css" rel="stylesheet">

  <style>
    #header-wrapper { margin: 0 auto; max-width: 1030px; }
    #content-wrapper { margin: 0 auto; max-width: 1000px; min-height: 467px; }
    #footer { background-color: #2a2730; margin-top: 30px; padding: 20px 15px; color: #99979c; text-align: right; }
    #footer-wrapper { margin: 0 auto; max-width: 1000px; }
    #footer-wrapper a { color:#fff;  }
    .markdown-body { padding: 30px; }
    @media (max-width: 767px) {
      .markdown-body { padding: 15px; }
    }
    .mono-font { font-family: "SFMono-Regular", Consolas, "Liberation Mono", Menlo, Courier; }
    .margin-right-10 { margin-right: 5px; }
    .gray { color: gray; }
    .small-font { font-size: 14px; }
  </style>
</head>
<body>
  <nav class="navbar navbar-default">
    <div class="container-fluid" id="header-wrapper">
      <div class="navbar-header">
        <a class="navbar-brand" href="/">李强的博客</a>
      </div>

      <div class="navbar-form navbar-right" role="search">
        <div class="form-group">
          <form action="/search">
            <input id="search-blog" name="q" type="text" class="form-control" placeholder="Search">
          </form>
        </div>
      </div>

    </div>
  </nav>

  <div class="panel panel-default" id="content-wrapper">
    <div class="panel-body markdown-body" id="content">
      {{range .Posts}}
      <p class="mono-font">
        <span class="margin-right-10 gray">{{.Date | format}}</span>
        <a href="/post/{{.Path}}">{{.Title}}</a>
      </p>
      {{end}}
    </div>
  </div>

  <div id="footer">
    <div id="footer-wrapper">
      © 2017-2018
      <a href="https://github.com/lqqyt2423">GitHub</a>
      <a href="mailto:974923609@qq.com">Email</a>
      <a href="/">李强的博客</a>
    </div>
  </div>

</body>
</html>
