import React from "react";

import "./App.css";

import { ThemeProvider, createTheme } from "@mui/material/styles";
import NavBar from "./components/navbar";
import Footer from "./components/footer";
import themeOptions from "./theme";
import "fontsource-roboto";
import { Routes } from "react-router-dom";
import { Route } from "react-router-dom";
import Home from "./pages/Home";
import History from "./pages/History";
import Config from "./pages/Config";
import Interests from "./pages/Interests";
import Login from "./pages/Login";
import Signup from "./pages/Signup";

// const beamsClient = new pusher.Client({
//     instanceId: "6e482588-a9a1-45a9-b786-2d367fc69eef"
// });

const theme = createTheme(themeOptions);

function App() {
    return (
        <>
            <ThemeProvider theme={theme}>
                <div className="App">
                    <NavBar />
                    <Routes>
                        <Route path="/" element={<Home />} />
                        <Route path="/History" element={<History />} />
                        <Route path="/Config" element={<Config />} />
                        <Route path="/Interests" element={<Interests />} />
                        <Route path="/Login" element={<Login />} />
                        <Route path="/Signup" element={<Signup />} />
                    </Routes>

                    <Footer></Footer>
                </div>
            </ThemeProvider>
        </>
    );
}

export default App;
