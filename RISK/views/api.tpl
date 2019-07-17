<body>
  <table>
    <thead>
      <tr>
        <th>Entity</th>
        <th>VaR</th>
      </tr>
    </thead>
    <tbody>
    
      {{range .temp}}
      <tr>
        <td>{{ .head }}</td>
        <td>{{ printf "%.2f" .vatr }}</td>
      </tr>
      {{ end }}
    </tbody>
  </table>
</body>