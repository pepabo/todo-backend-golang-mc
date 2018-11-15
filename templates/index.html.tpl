<!DOCTYPE html>
<html>
  <head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
    <link rel="stylesheet" href="https://fonts.googleapis.com/icon?family=Material+Icons">
    <link rel="stylesheet" href="https://code.getmdl.io/1.3.0/material.indigo-pink.min.css">
    <script defer src="https://code.getmdl.io/1.3.0/material.min.js"></script>
  </head>
  <body>
    <table class="mdl-data-table mdl-js-data-table mdl-shadow--2dp">
      <thead>
        <tr>
          <th class="mdl-data-table__cell--non-numeric">UA</th>
          <th>Created At</th>
        </tr>
      </thead>
      <tbody>
        {{range $i, $l := .Logs}}
        <tr>
          <td class="mdl-data-table__cell--non-numeric">{{ $l.UA }}</td><td>{{ $l.CreatedAt }}</td>
        </tr>
        {{end}}
      <tbody>
    </table>
  </body>
</html>
