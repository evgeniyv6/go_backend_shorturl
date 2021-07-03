import React, {Component} from 'react';

class Shorter extends React.Component {
  render() {
    return (
        <table className="table table-striped">
          <thead>
          <tr>
              <th>ID</th>
              <th>LINK</th>
              <th>STAT</th>
          </tr>
          </thead>
          {/*<div><strong><pre>{JSON.stringify(this.props.shorter_result_link, null, 2) }</pre></strong></div>*/}
          <tbody>
          {this.props.shorter_result_link &&
          <tr key={this.props.shorter_result_link.id}>
            <td>{this.props.shorter_result_link.id}</td>
            <td>{this.props.shorter_result_link.link}</td>
            <td>{this.props.shorter_result_link.stat}</td>
          </tr>
          }
          </tbody>
        </table>
    );
  }
}

class App extends Component {
  constructor(props) {
    super(props);
    this.state = {
      url: '',
      res:'',
      shorter_result_link: [],
    };

    this.create = this.create.bind(this);
    this.info = this.info.bind(this);
    this.handleChange = this.handleChange.bind(this);
  }

  componentDidMount() {

  }

  create(e) {
    e.preventDefault();
    console.log(JSON.stringify({
      url: this.state.url
    }))
    fetch("cut", {
      "method": "POST",
      "headers": {
        "content-type": "application/json",
        "accept": "application/json"
      },
      "body": JSON.stringify({
        url: this.state.url
      })
    })
        .then(response => response.json())
        .then(response => {
          this.setState({
            res: response
          })
        })
        .then(response => {
          console.log(response)
        })
        .catch(err => {
          console.log(err);
        });

  }

  info(e) {
    e.preventDefault();
    fetch(`${this.state.url}/info`, {
      "method": "GET",
      "headers": {
        "content-type": "application/json",
        "accept": "application/json"
      }
    })
        .then(response => response.json())
        .then(response => {
          this.setState({
            shorter_result_link: response
          })
        })
        .then(response => {
          console.log(`>>  ${JSON.stringify(this.state.shorter_result_link)}`)
        })
        .catch(err => {
          console.log(err);
        });
  }


  handleChange(changeObject) {
    this.setState(changeObject)
  }

  render() {
    return (
        <div className="container">
          <div className="row justify-content-center">
            <div className="col-md-8">
              <h1 className="display-4 text-center">Make a short link</h1>
              <form className="d-flex flex-column">
                <label htmlFor="name">
                  URL:
                  <input
                      name="url"
                      id="url"
                      type="text"
                      className="form-control"
                      value={this.state.url}
                      onChange={(e) => this.handleChange({ url: e.target.value })}
                      required
                  />
                </label>
                <button className="btn btn-primary" type='button' onClick={(e) => this.create(e)}>
                  Shorten
                </button>
                <br></br>
                <strong>{JSON.stringify(this.state.res.data)}</strong>
                <br></br>
                <button className="btn btn-info" type='button' onClick={(e) => {this.info(e)}}>
                  Get Stat
                </button>
                {/*<p>{JSON.stringify(this.state.shorter_result_link)}</p>*/}
                <br></br>
                <Shorter shorter_result_link={this.state.shorter_result_link.data} />
                <br></br>
              </form>
            </div>
          </div>
        </div>
    );
  }
}

export default App;