import React,{useState,useEffect} from 'react';
import AppBar from '@material-ui/core/AppBar';
import Button from '@material-ui/core/Button';
import Card from '@material-ui/core/Card';
import CardActions from '@material-ui/core/CardActions';
import CardContent from '@material-ui/core/CardContent';
import CardMedia from '@material-ui/core/CardMedia';
import CssBaseline from '@material-ui/core/CssBaseline';
import Grid from '@material-ui/core/Grid';
import Toolbar from '@material-ui/core/Toolbar';
import Typography from '@material-ui/core/Typography';
import { makeStyles } from '@material-ui/core/styles';
import Container from '@material-ui/core/Container';
import Link from '@material-ui/core/Link';
import ExitToAppIcon from '@material-ui/icons/ExitToApp';
import {useParams} from 'react-router-dom';
import { useHistory } from 'react-router-dom';
import CloudUploadIcon from '@material-ui/icons/CloudUpload';
import DeleteIcon from '@material-ui/icons/Delete';
import CloudDownloadIcon from '@material-ui/icons/CloudDownload';
import Alert from '@material-ui/lab/Alert';


function Copyright() {
  return (
    <Typography variant="body2" color="textSecondary" align="center">
      {'Copyright © '}
      <Link color="inherit" href="https://material-ui.com/">
        Your Website
      </Link>{' '}
      {new Date().getFullYear()}
      {'.'}
    </Typography>
  );
}

const useStyles = makeStyles((theme) => ({
  icon: {
    marginRight: theme.spacing(2),
  },
  heroContent: {
    backgroundColor: theme.palette.background.paper,
    padding: theme.spacing(8, 0, 6),
  },
  heroButtons: {
    marginTop: theme.spacing(4),
  },
  cardGrid: {
    paddingTop: theme.spacing(8),
    paddingBottom: theme.spacing(8),
  },
  card: {
    height: '100%',
    display: 'flex',
    flexDirection: 'column',
  },
  cardMedia: {
    paddingTop: '56.25%', // 16:9
  },
  cardContent: {
    flexGrow: 1,
  },
  footer: {
    backgroundColor: theme.palette.background.paper,
    padding: theme.spacing(6),
  },
}));

const cards = [1, 2, 3, 4, 5, 6, 7, 8, 9];

export default function Files() {
  const classes = useStyles();
  const [selectedFile, setSelectedFile] = useState({});
  const [isFilePicked, setIsFilePicked] = useState(false);
  const [currentPath,setCurrentPath] = useState(useParams().path);
  const [fileUploading,setFileUploading] = useState(false);
  const [fileUploaded,setFileUploaded] = useState(false);
  
  console.log(useParams().path);

  let history = useHistory();

  const changeHandler = (event) => {
    console.log(event.target.files[0]);
    setSelectedFile(event.target.files[0]);
    setIsFilePicked(true);
  };

  const handleSubmission = () => {
    const formData = new FormData();
    formData.append('file', selectedFile);
    formData.append('path',currentPath);
    setFileUploading(true);
    fetch(
      'http://localhost:8080/v1/MyCloud/upload',
      {
        method: 'POST',
        mode:"cors",
        credentials:"include",
        body: formData,
      }
    )
      .then((response) => response.json())
      .then((result) => {
        console.log('Success:', result);
        setFileUploading(false);
        setFileUploaded(true);
      })
      .catch((error) => {
        console.error('Error:', error);
      });
      setIsFilePicked(false);
  };


  const handleLogout = () => {
    fetch('http://localhost:8080/v1/MyCloud/logout',{
    method: "GET",
    mode:"cors",
    credentials:"include",
    headers: {"Content-type": "application/json; charset=UTF-8",/*"Origin":"http://localhost:3000"*/}
    })
    .then(response => response.json())
    .then(json => console.log(json))
    .catch(err => console.log('Request Failed',err))
    history.push("/")
  }  

  const listDirectory = async () => {
    const data = {path:currentPath}
    const response = await fetch('http://localhost:8080/v1/MyCloud/listDir',{
      method: "POST",
      mode:"cors",
      credentials:"include",
      body: JSON.stringify(data),
      headers: {"Content-type": "application/json; charset=UTF-8",/*"Origin":"http://localhost:3000"*/}
      });
    const jsonResponse = await response.json()
    if (response.status === 202) {
      console.log(jsonResponse)
    }else{
      console.log(jsonResponse)
    }
  }

  useEffect(() => {
    listDirectory();
  }, []);

  return (
    <React.Fragment>
      <CssBaseline />
      <AppBar position="relative">
        <Toolbar>
          <ExitToAppIcon />
          <Button onClick={handleLogout}>
            Logout
        </Button>
        </Toolbar>
      </AppBar>
      <main>
        {/* Hero unit */}
        <div className={classes.heroContent}>
          <Container maxWidth="sm">
            <Typography component="h1" variant="h2" align="center" color="textPrimary" gutterBottom>
              My Files
            </Typography>
            <Typography variant="h5" align="center" color="textSecondary" paragraph>
              Current Directory:{currentPath}
            </Typography>
            <div className={classes.heroButtons}>
              <Grid container spacing={2} justify="center">
              <Grid item>
                <Button
                    variant="contained"
                    component="label"
                    color="primary"
                    className={classes.button}
                  >Choose a file to upload
                <input type="file" name="file" hidden onChange={changeHandler} />
                </Button>
                </Grid>

                <Grid item>
                  <Button variant="outlined" color="primary">
                    Create a directory
                  </Button>
                </Grid>
              </Grid>

              {isFilePicked ? (
                  <React.Fragment>
                    <Grid item>
                      <Typography variant="h6" gutterBottom>
                        <br></br>Filename: {selectedFile.name}
                      </Typography>
                    </Grid>
                    <Grid item>
                      <Typography variant="h6" gutterBottom>
                        Size in bytes: {selectedFile.size}
                      </Typography>
                    </Grid>
                    <Grid item> 
                      <Button
                      variant="contained"
                      color="primary"
                      className={classes.button}
                      startIcon={<CloudUploadIcon />}
                      onClick={handleSubmission}
                      >
                      Upload
                      </Button>
                    </Grid>
                  </React.Fragment>

                ) : (
                  <p></p>
                )}

            </div>
            {fileUploading && <Alert severity="info">File Uploading</Alert>}
            {fileUploaded && <Alert onClose={() => setFileUploaded(false)} severity="success">File Uploaded</Alert>}
          </Container>
        </div>
        <Container className={classes.cardGrid} maxWidth="md">
          {/* End hero unit */}
          <Grid container spacing={4}>
            {cards.map((card) => (
              <Grid item key={card} xs={12} sm={6} md={4}>
                <Card className={classes.card}>
                  <CardMedia
                    className={classes.cardMedia}
                    image="https://source.unsplash.com/random"
                    title="Image title"
                  />
                  <CardContent className={classes.cardContent}>
                    <Typography gutterBottom variant="h5" component="h2">
                      Heading
                    </Typography>
                    <Typography>
                      This is a media card. You can use this section to describe the content.
                    </Typography>
                  </CardContent>
                  <CardActions>
                    <Button
                      variant="contained"
                      color="primary"
                      className={classes.button}
                      startIcon={<CloudDownloadIcon />}
                    >
                      Download
                    </Button>
                    <Button
                      variant="contained"
                      color="secondary"
                      className={classes.button}
                      startIcon={<DeleteIcon />}
                    >
                      Delete
                    </Button>
                  </CardActions>
                </Card>
              </Grid>
            ))}
          </Grid>
        </Container>
      </main>
      {/* Footer */}
      <footer className={classes.footer}>
        <Typography variant="h6" align="center" gutterBottom>
          Footer
        </Typography>
        <Typography variant="subtitle1" align="center" color="textSecondary" component="p">
          Something here to give the footer a purpose!
        </Typography>
        <Copyright />
      </footer>
      {/* End footer */}
    </React.Fragment>
  );
}