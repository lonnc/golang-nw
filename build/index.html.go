package build

const index = `
<!DOCTYPE html>
<html>
<head>
<title>{{ .Name }}</title>
<style>
* {
  margin: 0;
  padding: 0;
}

html, body {
  height: 100%;
}

</style>
</head>
<body>
  <pre id="state"></pre>
  <script src="script.js"></script>
</body>
</html>
`
