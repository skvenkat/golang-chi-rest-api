import "./App.css";
import {Route, BrowserRouter, Routes, Link} from "react-router-dom";
import {HomePage} from "./HomePage";
import React from "react";
import {AboutPage} from "./AboutPage";
function App() {
  return (<div className="App">
    <BrowserRouter>
        <Navbar/>
        <Routes>
          <Route path="/" element={<HomePage/>} />
          <Route path="/webapp/about" element={<AboutPage/>} />
        </Routes>
    </BrowserRouter>
  </div>)
}

export default App;

function Navbar() {
  return (
    <nav className="flex items-center justify-between flex-wrap bg-indigo-500 p-6">
      <div className="flex items-center text-white">
        <span className="font-semibold text-xl">React-Typescript Example</span>
      </div>
      <div >
        <div className="text-sm">
          <Link
            to={"/"}
            className="mt-4 lg:inline-block lg:mt-0 text-indigo-200 hover:text-white mr-4"
          >
            Home
          </Link>
          <Link
            to={"/webapp/about"}
            className="mt-4 lg:inline-block lg:mt-0 text-indigo-200 hover:text-white mr-4"
          >
            About
          </Link>
        </div>
      </div>
    </nav>
  )
}
