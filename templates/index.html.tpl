<!DOCTYPE html>
<html>
  <head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
  </head>
  <body>
    <h1>{{.Title}}</h1>
    <table>
      <thead>
        <tr>
          <th>Title</th>
          <th>Created At</th>
        </tr>
      </thead>
      <tbody>
        {{range $i, $m := .Memos}}
        <tr>
          <td>{{ $m.Title }}</td><td>{{ $m.CreatedAt }}</td>
        </tr>
        {{end}}
      <tbody>
    </table>
  </body>
</html>
