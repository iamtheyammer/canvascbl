import React from 'react';
import { Redirect } from 'react-router-dom';
import './index.css';

function Home(props) {
  return (
    <>
      <Redirect to="/dashboard" />
    </>
  );
}

// const ConnectedHome = connect((state) => ({
//   getSessionId: state.components.home.getSessionId,
//   signInButtonAvailability: state.components.home.signInButtonAvailability,
//   loading: state.loading,
//   error: state.error,
//   session: state.plus.session
// }))(Home);

export default Home;
