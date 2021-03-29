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
import TextField from '@material-ui/core/TextField';
import Breadcrumbs from '@material-ui/core/Breadcrumbs';


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
    height: 0,
  },
  cardContent: {
    flexGrow: 1,
  },
  footer: {
    backgroundColor: theme.palette.background.paper,
    padding: theme.spacing(6),
  },
}));


export default function Files() {
  const classes = useStyles();
  const [selectedFile, setSelectedFile] = useState({});
  const [isFilePicked, setIsFilePicked] = useState(false);
  const [currentPath,setCurrentPath] = useState(useParams().path);
  const [fileUploading,setFileUploading] = useState(false);
  const [fileUploaded,setFileUploaded] = useState(false);
  const [createDir,setCreateDir] = useState(false);
  const [newDirInfo,setNewDirInfo] = useState({dirName:'',path:useParams().path});
  const [errorMessage,setErrorMessage] = useState('');
  const [directoryEntries,setDirectoryEntries] = useState([]);
  const [deleteWarning,setDeleteWarning] = useState({show:false,entryToDeleteName:''});
  const [dirSize,setDirSize] = useState(0);
  
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
        setDirectoryEntries([...directoryEntries,result.newFile])
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
      console.log(jsonResponse.dirEntries)
      setDirectoryEntries(jsonResponse.dirEntries)
    }else{
      console.log(jsonResponse.error)
    }
  }

  const handleChange = (e) => {
    const name = e.target.name;
    const value = e.target.value;
    setNewDirInfo({ ...newDirInfo, [name]: value });
  };  

  const handleNewDirSubmit = (e) => {
    e.preventDefault();
    console.log(newDirInfo);
    fetch('http://localhost:8080/v1/MyCloud/createDir', {
      method: "POST",
      mode:"cors",
      credentials:"include",
      body: JSON.stringify(newDirInfo),
      headers: {"Content-type": "application/json; charset=UTF-8",/*"Origin":"http://localhost:3000"*/}
    })
    .then(response => response.json())
    .then((json) => {
      console.log(json);
      if(json.error){
        setErrorMessage(json.error);
      }else{
        setCreateDir(false)
        setDirectoryEntries([...directoryEntries,json.newDirectory])
      }
    });
  }

  function handleChangeDirectory(dirName){
    var str = "_";
    var dir = str.concat(dirName);
    var newDir = currentPath.concat(dir);
    setCurrentPath(newDir);
    history.push(`${newDir}`);
  }

  function handleBreadcrumb(event) {
    event.preventDefault();
    var targetDir = event.target.text;
    var n = currentPath.indexOf(targetDir);
    var targetPath = currentPath.substr(0, n+targetDir.length)
    setCurrentPath(targetPath);
    history.push(`${targetPath}`);
  }

  function handleDeleteEntry(entryObject){
    console.log("DELETING ",entryObject);
    const entryName = entryObject.name;
    const isDir = entryObject.isDirectory;
    const data = {name:entryName,isDirectory:isDir,path:currentPath}
    console.log("DATA IS",JSON.stringify(data))
    fetch('http://localhost:8080/v1/MyCloud/delete', {
      method: "POST",
      mode:"cors",
      credentials:"include",
      body: JSON.stringify(data),
      headers: {"Content-type": "application/json; charset=UTF-8",}
    })
    .then(response => response.json())
    .then((json) => {
      console.log(json);
      if(json.error){
        setErrorMessage(json.error);
      }else{
        setDeleteWarning({show:false,entryToDeleteName:''});
        let newEntries = directoryEntries.filter((entry) => entry.name !== entryObject.name);
        setDirectoryEntries(newEntries);
      }
    }); 
  }

  const getDirectorySize = () => {
    fetch('http://localhost:8080/v1/MyCloud/size',{
    method: "GET",
    mode:"cors",
    credentials:"include",
    headers: {"Content-type": "application/json; charset=UTF-8",/*"Origin":"http://localhost:3000"*/}
    })
    .then(response => response.json())
    .then(json => setDirSize(json.totalSize))
    .catch(err => console.log('Request Failed',err))
    
  }  

  useEffect(() => {
    listDirectory();
    getDirectorySize();
  }, [currentPath]);

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
              Used space:{dirSize} GB / 2GB
              <Breadcrumbs aria-label="breadcrumb" align="inherit">
              {currentPath.split("_").map((dirEntry) => (
                      <Link key={dirEntry} color="inherit" href="/" align="inherit" onClick={handleBreadcrumb}>
                      {dirEntry}
                    </Link>
              ))}
              </Breadcrumbs>
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
                  <Button variant="outlined" color="primary" onClick={() => setCreateDir(true)}>
                    Create a directory
                  </Button>
                </Grid>
                {createDir &&
                      <div className={classes.paper}>
                        <br></br>
                      <form className={classes.form} noValidate>
                        <Grid container spacing={2}>
                          <Grid item xs={12}>
                            <TextField
                              autoComplete="fname"
                              name="dirName"
                              variant="outlined"
                              required
                              fullWidth
                              id="dirName"
                              label="Directory Name"
                              onChange={handleChange}
                              autoFocus
                            />
                          </Grid>
                        </Grid>
                        <Button
                          type="submit"
                          fullWidth
                          variant="contained"
                          color="primary"
                          className={classes.submit}
                          onClick={handleNewDirSubmit}
                        >
                          Create Directory
                        </Button>
                        <Button variant="contained" color="secondary" onClick={() => setCreateDir(false)}>
                          Cancel
                        </Button>
                      </form>
                      {errorMessage && <Alert onClose={() => setErrorMessage("")} severity="error">{errorMessage}</Alert>}
                    </div>
                }
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
                        Size: {(selectedFile.size/ 1024.0 / 1024.0).toFixed(3)} MB
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
                      <Button
                      variant="contained"
                      color="secondary"
                      className={classes.button}
                      onClick={() => setIsFilePicked(false)}
                      >
                      Cancel
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
            {directoryEntries.map((dirEntry) => (
              <Grid item key={dirEntry.name} xs={12} sm={6} md={4}>
                <Card className={classes.card}>
                  <CardMedia
                    className={classes.cardMedia}
                    image = {dirEntry.icon}
                    title="Image title"
                  />
                  <CardContent className={classes.cardContent}>
                    <Typography gutterBottom variant="h6" component="h2">
                      {dirEntry.name}
                    </Typography>
                  </CardContent>
                  <CardActions>
                    {!dirEntry.isDirectory ?(
                    <Button
                      variant="contained"
                      color="primary"
                      href={dirEntry.link}
                      className={classes.button}
                      startIcon={<CloudDownloadIcon />}
                    >
                      Open
                    </Button>
                    ): (
                      <Button
                      variant="contained"
                      color="primary"
                      href="#"
                      onClick={() => handleChangeDirectory(dirEntry.name)}
                      className={classes.button}
                      startIcon={<CloudDownloadIcon />}
                    >
                      Open
                    </Button>
                    )}
                    <Button
                      variant="contained"
                      color="secondary"
                      className={classes.button}
                      startIcon={<DeleteIcon />}
                      onClick={() => setDeleteWarning({show:true,entryToDeleteName:dirEntry.name})}
                    >
                      Delete
                    </Button>
                  </CardActions>
                  {deleteWarning.show && deleteWarning.entryToDeleteName === dirEntry.name  &&
                  <Alert
                    severity="warning"
                    action={
                    <div>
                      <Button color="inherit" size="small" onClick={() => handleDeleteEntry(dirEntry)}>
                        DELETE
                      </Button>
                      <Button color="inherit" size="small" onClick={() => {setDeleteWarning({show:false,entryToDeleteName:""})}}>
                        CANCEL
                      </Button>
                    </div>
                    }
                  >
                    Are you sure you want to delete it?
                  </Alert>
                  }
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