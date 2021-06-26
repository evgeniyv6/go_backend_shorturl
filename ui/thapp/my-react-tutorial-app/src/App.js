import React, {Component} from 'react';

class Friends extends Component {
  render() {
    return (
        <table className="table table-striped">
          <thead>
          <tr>
            <th>ID</th>
            <th>Original URL</th>
            <th>Call statistics</th>
          </tr>
          </thead>
          <tbody>
          {this.props.friends && this.props.friends.length>0 && this.props.friends.map(friend => {
            return <tr key={friend.id}>
              <td>{friend.id}</td>
              <td>{friend.link}</td>
              <td>{friend.stat}</td>
            </tr>
          })}
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
      friends: [],
    };

    this.create = this.create.bind(this);
    this.info = this.info.bind(this);
    this.handleChange = this.handleChange.bind(this);
  }

  componentDidMount() {
    fetch("http://0.0.0.0:8079/kvlzpsVUPH/info", {
      "method": "GET",
      "headers": {
        "content-type": "application/json",
        "accept": "application/json"
      }
    })
        .then(response => response.json())
        .then(response => {
          this.setState({
            friends: response
          })
        })
        .then(response => {
          console.log(`resp_one = ${response}`)})
        .catch(err => { console.log(err);
        });
  }

  info(e) {
    e.preventDefault();

    console.log(`url-->> ${this.state.url}`)

    fetch(this.state.url, {
      "method": "GET",
      "headers": {
        "content-type": "application/json",
        "accept": "application/json"
      }
    })
        .then(response => response.json())
        .then(response => {
          console.log(response);
        })
        .catch(err => {
          console.log(err);
        });
  }


  handleChange(event) {
    this.setState(event)
  }

  handleSubmit(event) {
    event.preventDefault();
    fetch("http://0.0.0.0:8079/kvlzpsVUPH/info", {
      "method": "GET",
      "headers": {
        "content-type": "application/json",
        "accept": "application/json"
      }
    })
        .then(response => response.json())
        .then(response => {
          console.log(`response123 - ${response}`)
        })
        .then(url => this.setState(url));
  }

// async postData(url = '', data = {}) {
//   // Default options are marked with *
//   const response = await fetch(url, {
//     method: 'GET', // *GET, POST, PUT, DELETE, etc.
//     mode: 'no-cors', // no-cors, *cors, same-origin
//     cache: 'no-cache', // *default, no-cache, reload, force-cache, only-if-cached
//     credentials: 'same-origin', // include, *same-origin, omit
//     headers: {
//       'Content-Type': 'application/json'
//       // 'Content-Type': 'application/x-www-form-urlencoded',
//     },
//     redirect: 'follow', // manual, *follow, error
//     referrerPolicy: 'no-referrer' // no-referrer, *client
//     //body: JSON.stringify(data) // body data type must match "Content-Type" header
//   });
//   return await response; // parses JSON response into native JavaScript objects
// }

  create(e) {
    // add entity - POST
    e.preventDefault();


    console.log(JSON.stringify({
      url: this.state.link_url
    }))
    fetch("http://0.0.0.0:8079/cut", {
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
          console.log(response)
        })
        .catch(err => {
          console.log(err);
        });

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
                <button className="btn btn-info" type='button' onClick={(e) => {
                  this.info(e);

                }}>
                  Get Stat
                </button>
              </form>
              <Friends friends={this.state.friends} />
            </div>
          </div>
        </div>

    );
  }
}

export default App;