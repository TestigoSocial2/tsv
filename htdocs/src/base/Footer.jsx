import React from 'react';

class Footer extends React.Component {
  render() {
    return(
      <div id="bottom" className="row">
        <div className="inner-row">
          <div className="col-md-3">
            <p>Una iniciativa de:</p>
            <a href="http://www.tm.org.mx" target="_blank" className="tm"></a>
          </div>
          <div className="col-md-3">
            <p>Con el apoyo de:</p>
            <a href="https://www.gov.uk/government/world/organisations/british-embassy-mexico-city.es-419" target="_blank" className="uk"></a>
          </div>
          <div className="col-md-6">
            <div className="social">
              <a href="https://www.twitter.com/testigo_social" className="tw"></a>
              <a href="https://www.facebook.com/testigosocial" className="fb"></a>
            </div>
          </div>
        </div>
      </div>
    );
  }
}

export default Footer;
