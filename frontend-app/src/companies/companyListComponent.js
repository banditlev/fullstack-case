import {useEffect, useState} from "react";
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableContainer from '@material-ui/core/TableContainer';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import DeleteIcon from '@material-ui/icons/Delete';
import Paper from '@material-ui/core/Paper';
import {
  Button,
  Card,
  CardContent,
  IconButton,
  TextField
} from "@material-ui/core";

function CompanyListComponent() {
  const [error, setError] = useState(null);
  const [isLoaded, setIsLoaded] = useState(false);
  const [companies, setCompanies] = useState([]);
  const [count, setCount] = useState(0);
  const [inputs, setInputs] = useState({CVR: "", Name: "", Address: ""});

  useEffect(() => {
    fetchCompanyList()
  }, [])

  const handleInputChange = (event) => {
    event.persist();
    setInputs(inputs => ({...inputs, [event.target.name]: event.target.value}));
  }

  const fetchCompanyList = () => {
    fetch("http://localhost:8080/companies")
    .then(res => res.json())
    .then(
        (result) => {
          setIsLoaded(true);
          setCompanies(result?.Companies);
          setCount(result?.Count)
        },
        (error) => {
          setIsLoaded(true);
          setError(error);
        }
    )
  }

  const createCompany = (event) => {
    if (event) {
      event.preventDefault();
      const requestOptions = {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify(
            {CVR: +inputs.CVR, Name: inputs.Name, Address: inputs.Address})
      };
      fetch('http://localhost:8080/companies', requestOptions)
      .then(() => {
        resetForm();
        fetchCompanyList();
      });
    }
  }

  const deleteCompany = (CVR) => {
    if (CVR !== undefined) {
      const requestOptions = {
        method: 'DELETE',
      };
      fetch('http://localhost:8080/companies/' + CVR, requestOptions)
      .then(() => fetchCompanyList());
    }
  }

  const resetForm = () => {
    Array.from(document.querySelectorAll("input")).forEach(
        input => (input.value = "")
    );
    setInputs(() => ({cvr: "", name: "", address: ""}));
  };

  if (error) {
    return <div>Error: {error.message}</div>;
  } else if (!isLoaded) {
    return <div>Loading...</div>;
  } else {
    return (
        <div>
          <Card>
            <CardContent>
              <form onSubmit={createCompany} className="company-input"
                    autoComplete="off">
                <TextField onChange={handleInputChange} name="CVR"
                           required
                           value={inputs.CVR}
                           style={{marginRight: 20}}
                           label="CVR"/>
                <TextField onChange={handleInputChange} name="Name"
                           required
                           value={inputs.Name}
                           style={{marginRight: 20}}
                           label="Name"/>
                <TextField onChange={handleInputChange} name="Address"
                           required
                           value={inputs.Address}
                           style={{marginRight: 20}} label="Address"/>
                <Button style={{marginTop: 15}} type="submit"
                        variant="contained" color="primary">
                  Add
                </Button>
              </form>
            </CardContent>
          </Card>
          <div>
            <h3> Total companies {count} </h3>
            <TableContainer component={Paper}>
              <Table aria-label="simple table">
                <TableHead>
                  <TableRow>
                    <TableCell>CVR</TableCell>
                    <TableCell align="right">Name</TableCell>
                    <TableCell align="right">Address</TableCell>
                    <TableCell align="right"> </TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                  {companies?.map((company, index) => (
                      <TableRow key={index}>
                        <TableCell component="th" scope="row">
                          {company.CVR}
                        </TableCell>
                        <TableCell align="right">{company.Name}</TableCell>
                        <TableCell align="right">{company.Address}</TableCell>
                        <TableCell align="right">
                          <IconButton onClick={() => deleteCompany(company.CVR)}
                                      aria-label="delete">
                            <DeleteIcon/>
                          </IconButton>
                        </TableCell>
                      </TableRow>
                  ))}
                </TableBody>
              </Table>
            </TableContainer>
          </div>
        </div>
    );
  }
}

export default CompanyListComponent;
