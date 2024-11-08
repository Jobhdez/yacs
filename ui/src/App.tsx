import React from "react";
import { Grid, Box } from "@mui/material";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";

import Compiler from "./Components/Compiler";

function App() {
  return (
    <Router>
      <Routes>
        <Route exact path="/" element={<Compiler />} />
      </Routes>
    </Router>
  );
}

export default App;
