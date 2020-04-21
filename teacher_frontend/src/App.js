import React from "react";
import { HashRouter as Router, Switch, Route, Link } from "react-router-dom";
import * as ReactGA from "react-ga";
import "./App.css";
import env from "./util/env";
import Home from "./components/Home";
import Dashboard from "./components/Dashboard";

function App(props) {
  ReactGA.initialize(env.googleAnalyticsId);

  ReactGA.pageview("/");

  return (
    <Router>
      <Switch>
        <Route exact path="/" component={Home} />
        <Route path="/dashboard" component={Dashboard} />
        <Route
          status={404}
          render={() => (
            <div align="center">
              <p style={{ color: "white" }}>404 Not Found</p>
              <Link to="/">Home</Link>
            </div>
          )}
        />
      </Switch>
    </Router>
  );
}

export default App;
