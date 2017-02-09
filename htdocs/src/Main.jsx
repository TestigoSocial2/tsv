import React from 'react';
import Menu from './base/Menu.jsx';
import Footer from './base/Footer.jsx';
import PrivacyNotice from './base/PrivacyNotice.jsx';

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

        {/* Privacy Notice */}
        <PrivacyNotice />
      </div>
    );
  }
}

export default Main;
