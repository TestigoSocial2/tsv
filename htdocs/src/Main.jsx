import React from 'react';
import Menu from './base/Menu.jsx';
import Footer from './base/Footer.jsx';

class Main extends React.Component {
  componentDidMount() {
    $('[data-toggle="tooltip"]').tooltip();
  }

  render() {
    return(
      <div className="container-fluid">
        {/* Top Menu */}
        <Menu />

        {/* Routes */}
        <div id="route">
          {this.props.children}
        </div>

        {/* Footer */}
        <Footer />
      </div>
    );
  }
}

export default Main;
