<!DOCTYPE html>
<html>
  <head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
  </head>
  <body>
    <h1>{{.Title}}</h1>
    <form method="post">
      <input type="text" name="title" />
      <input type="text" name="body" />
      <input type="submit"/>
    </form>
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
