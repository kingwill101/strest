const express = require("express")

const app = express();

app.get('/user', (req, res) => {
  res.send({
    username: req.query.name,
    id: 123,
    premium: false,
  })
})

app.listen(3001)

