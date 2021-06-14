import logo from './logo.svg';
import './App.css';
import React from "react";
import {
  BrowserRouter as Router,
  Switch,
  Route,
  Link
} from "react-router-dom";
// function App() {
//   return (
//     <div className="App">
//       <h1>Implicit Grant Type!</h1>
//       <div><a href="http://192.168.100.101:8080/auth/realms/learningApp/protocol/openid-connect/auth?client_id=client-react-implicit&response_type=token">Login</a></div>
//       <div><a>Service</a></div>
      
//     </div>
//   );
// }

function App() {
  return (
    <Router>
      <div>
        <div className="App">
          <h1>Implicit Grant Type!</h1>
        </div>    
        <nav>
          <ul>
            <li>
              <Link to="/">Home</Link>
            </li>
            <li>
              <Link to="/login">Login</Link>
            </li>
            <li>
              <Link to="/callback">Callback</Link>
            </li>
            <li>
              <Link to="/service">Service</Link>
            </li>
            <li>
              <Link to="/logout">Logout</Link>
            </li>
          </ul>
        </nav>

        {/* A <Switch> looks through its children <Route>s and
            renders the first one that matches the current URL. */}
        <Switch>
          <Route path="/login">
            <Login />
          </Route>
          <Route path="/callback">
            <Callback />
          </Route>

          <Route path="/service">
            <Service />
          </Route>
          <Route path="/logout">
            <Logout />
          </Route>
          <Route path="/">
            <Home />
          </Route>
        </Switch>
      </div>
    </Router>
  );
}

function Home() {
  return <h2>Home</h2>;
}

function Login() {
  window.location = "http://192.168.100.101:8080/auth/realms/learningApp/protocol/openid-connect/auth?client_id=client-react-implicit&response_type=token&redirect_uri=http://localhost:3000/callback"
  return null;
}
function Callback() {
  //Pegar o access_token
  return <h2>Callback</h2>;
}

function Service() {
  return <h2>Users</h2>;
}

function Logout() {
  return <h2>Users</h2>;
}

export default App;
