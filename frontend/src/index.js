import React,{useState} from 'react';
import ReactDom from 'react-dom';
import {BrowserRouter as Router,Route,Switch} from 'react-router-dom';
import {
  teal,
  lightBlue,
  deepOrange
} from "@material-ui/core/colors";
import { createMuiTheme, ThemeProvider } from "@material-ui/core/styles";

import SignUp from "./SignUp"
import SignIn from "./SignIn"
import HomePage from "./HomePage"
import Files from "./Files"


const App = () => {

  const [darkState, setDarkState] = useState(true);
  const palletType = darkState ? "dark" : "light";
  const mainPrimaryColor = darkState ? teal[500] : lightBlue[500];
  const mainSecondaryColor = darkState ? deepOrange[900] : lightBlue[500];
  const darkTheme = createMuiTheme({
    palette: {
      type: palletType,
      primary: {
        main: mainPrimaryColor
      },
      secondary: {
        main: mainSecondaryColor
      }
    }
  });


  return (
    <ThemeProvider theme={darkTheme}>
    <Router>
      <Switch>
        <Route exact path='/'>
          <HomePage/>
        </Route>
        <Route exact path='/signup'>
          <SignUp/>
        </Route>
        <Route exact path='/signin'>
          <SignIn/>
        </Route>
        <Route exact path='/files/:path' children={<Files/>}></Route>        
      </Switch>
    </Router>
    </ThemeProvider>
  )
}

ReactDom.render(<App/>,document.getElementById('root'));