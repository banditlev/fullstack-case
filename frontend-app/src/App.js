import './App.css';
import CompanyListComponent from './companies/companyListComponent'

function App() {
  return (
      <div className="App">
        <header className="App-header">
          <h1>
            List of companies
          </h1>
        </header>
        <div className="list-container">
          <CompanyListComponent/>
        </div>
      </div>
  );
}

export default App;
