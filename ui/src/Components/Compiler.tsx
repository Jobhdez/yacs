import React, { useState } from "react";
import {
  Box,
  Button,
  Typography,
  FormControl,
  Grid,
  InputLabel,
  MenuItem,
  Select,
} from "@mui/material";
import Editor from "@monaco-editor/react";

function Compiler() {
  const [ExpCode, setExpCode] = useState("");
  const [CompiledCode, setCompiledCode] = useState("");
  const [example, setExample] = useState("");

  const handleChange = (event) => {
    setExample(event.target.value);
    setExpCode(event.target.value);
  };

  const IfHandle = () => {
    const exampleCode = "(let ((x 3)) (if (< x 4) 3 4))";
    setExample(exampleCode);
    setExpCode(exampleCode);
  };

  const whileHandle = () => {
    const exampleCode =
      "(let ((i 0)) (while (< i 5) (begin i (set i (+ i 1)))))";
    setExample(exampleCode);
    setExpCode(exampleCode);
  };

  const while2Handle = () => {
    const exampleCode =
      "(let ((sum 0)) (let ((i 0)) (begin (while (< i 5) (begin (set sum (+ sum 3)) (set i (+ i 1)))) sum)))";
    setExample(exampleCode);
    setExpCode(exampleCode);
  };

  const varHandle = () => {
    const exampleCode = "(let ((x 2)) x)";
    setExample(exampleCode);
    setExpCode(exampleCode);
  };

  const lispApi = "http://localhost:1234/api/compiler";

  function Compile() {
    const data = JSON.stringify({ exp: ExpCode });

    fetch(lispApi, {
      method: "POST",
      headers: {
        Accept: "application/json",
        "Content-Type": "application/json",
      },
      body: data,
      mode: "cors",
    })
      .then((response) => response.json())
      .then((data) => {
        setCompiledCode(data.exp);
      });
  }

  return (
    <Grid container spacing={2}>
      <Grid container spacing={2}>
        <Grid item xs={6}>
          {" "}
          <FormControl size="small">
            <InputLabel
              sx={{
                color: "black",
                fontFamily: "monospace",
              }}
            >
              Programs
            </InputLabel>
            <Select
              label="Examples"
              value={example}
              onChange={handleChange}
              sx={{
                width: 250,
                height: 50,
                backgroundColor: "#3CB371",
              }}
            >
              <MenuItem
                value="(let ((x 3)) (if (< x 4) 3 4))"
                onClick={IfHandle}
              >
                If statement
              </MenuItem>
              <MenuItem
                value="(let ((i 0)) (while (< i 5) (begin i (set i (+ i 1)))))"
                onClick={whileHandle}
              >
                while statement
              </MenuItem>
              <MenuItem
                value="(let ((sum 0)) (let ((i 0)) (begin (while (< i 5) (begin (set sum (+ sum 3)) (set i (+ i 1)))) sum)))"
                onClick={while2Handle}
              >
                while 2nd eg
              </MenuItem>
              <MenuItem value="(let ((x 2)) x)" onClick={varHandle}>
                Variable
              </MenuItem>
            </Select>
          </FormControl>
        </Grid>
      </Grid>
      <Grid container spacing={6}>
        <Grid item xs={8}>
          <Box sx={{ height: "80vh" }}>
            <Editor
              theme="vs-dark"
              width="100%"
              height="100%"
              language="scheme"
              value={ExpCode} // Use value prop instead of defaultValue
              onChange={(val) => {
                setExpCode(val || "");
              }}
            />
          </Box>
        </Grid>
        <Grid item xs={4}>
          <Box
            sx={{
              height: "80vh",
              width: "150%",
              overflowY: "auto",
              backgroundColor: "#f5f5f5",
              padding: "10px",
              borderRadius: "4px",
            }}
          >
            <Typography variant="h6">x86 program:</Typography>
            <pre>
              <code>{CompiledCode}</code>
            </pre>
          </Box>
        </Grid>
      </Grid>

      <Button
        onClick={Compile}
        variant="contained"
        sx={{
          backgroundColor: "#3CB371",
          color: "black",
          "&:hover": {
            backgroundColor: "darkgray",
          },
          marginTop: "20px",
        }}
      >
        Compile Exp
      </Button>
    </Grid>
  );
}

export default Compiler;
