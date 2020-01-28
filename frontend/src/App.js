import React from 'react';
import { HashRouter as Router, Route, Switch, Link } from 'react-router-dom';
import * as ReactGA from 'react-ga';
import { isMobile } from 'react-device-detect';

import ConnectedHome from './components/Home';
import ConnectedOAuth2Response from './components/OAuth2Response';

import Dashboard from './components/Dashboard';
import env from './util/env';

function App(props) {
  ReactGA.initialize(env.googleAnalyticsId);

  ReactGA.pageview('/');

  if (!isMobile) {
    document.body.classList.add('background');
  }

  return (
    <Router>
      <Switch>
        <Route exact path="/" component={ConnectedHome} />
        <Route
          exact
          path="/oauth2response"
          component={ConnectedOAuth2Response}
        />
        <Route path="/dashboard" component={Dashboard} />
        <Route
          status={404}
          render={() => (
            <div align="center">
              <p color="#fffff">404 Not Found</p>
              <Link to="/">Home</Link>
            </div>
          )}
        />
      </Switch>
    </Router>
  );
}

export default App;
