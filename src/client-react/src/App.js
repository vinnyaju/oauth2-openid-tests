import './App.css';
import React from "react";
import {
  BrowserRouter as Router,
  Switch,
  Route,
  Link
} from "react-router-dom";
class App extends React.Component {
  constructor(props) {
    super(props)
    this.state = {
      access_token : "",
      expires_in: "",
      session_state: "",
      token_type: ""
    }
  }

  setStateValue = (k, v) => {
    if (this.state[k] !== v)
      this.setState({[k]: v})
  }

  onCheckStateClick = () => {
    console.log(this.state)
  }

  render() {
    return (
      <Router>
        <div>
          <div className="App">
            <h1>Implicit Grant Type!</h1>
          </div>
          <button onClick={this.onCheckStateClick}>Check stete.</button>
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
              <Callback setStateValue = {this.setStateValue} />
            </Route>

            <Route path="/service">
              <Service accessToken = {this.state.access_token} />
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
}

function Home() {
  return <h2>Home</h2>;
}

function Login() {
  window.location = "http://192.168.100.101:8080/auth/realms/learningApp/protocol/openid-connect/auth?client_id=client-react-implicit&response_type=token&redirect_uri=http://localhost:3000/callback&scope=getUtilitiesService"
  return null;
}
function Callback(props) {
  //Pegar o access_token
  const hashStr = window.location.hash
  const hashMap = hashStr.substr(1).split("&").reduce((accum, item) => {
      //Adicionar item ao acumulador (acc)
      const kv = item.split("=")
      accum[kv[0]] = kv[1]
      //Retornar o Acumulador.
      return accum
    }, {})
  //Guardando o State...
  const {setStateValue} = props

  setStateValue("access_token", hashMap.access_token)
  setStateValue("expires_in", hashMap.expires_in)
  setStateValue("session_state", hashMap.session_state)
  setStateValue("token_type", hashMap.token_type)

  return <h2>Callback</h2>;
}

class Service extends React.Component{
  constructor (props) {
    super(props)
    this.state = {
      data: {}
    }
  }

  componentDidMount() {
    const {accessToken} = this.props

    //Fazendo um POST com Form para acessar o recurso
    const formData = new FormData()
    formData.append("access_token", accessToken)
    fetch("http://localhost:4000/billing/v1/utilities", {
      method: "POST",
      body: formData
    })
    .then(response => response.json())
    .then(data => {
      console.log(data)
      this.setState({data})
    })

  }

  render() {
    return <div>
      <h2>Service</h2>
      <div>{JSON.stringify(this.state.data)}</div>
    </div>
  }
}

function Logout() {
  return <h2>Users</h2>;
}

export default App;
