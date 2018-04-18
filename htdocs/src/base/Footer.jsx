import React from 'react';

class Footer extends React.Component {
  render() {
    return(
      <div id="bottom" className="row">
        <div className="inner-row hide-row">
          <div className="col-md-6">
            <div className="logo">
              <p>Una iniciativa de:</p>
              <a href="http://www.tm.org.mx" target="_blank" className="tm"></a>
            </div>
            <div className="logo">
              <p>Con el apoyo de:</p>
              <a href="https://www.gov.uk/government/world/organisations/british-embassy-mexico-city.es-419" target="_blank" className="uk"></a>
            </div>
          </div>
          <div className="col-md-6 links">
            <div className="social">
              <a href="https://www.twitter.com/testigo_social" className="tw"></a>
              <a href="https://www.facebook.com/testigosocial" className="fb"></a>
              <div className="clear"></div>
            </div>
            <span data-toggle="modal" data-target="#privacyNotice">Aviso de Privacidad</span>
          </div>
        </div>
        <div className="inner-row mobile-row">
          <div className="col-md-6 links">
            <div className="social">
              <a href="https://www.facebook.com/testigosocial" target="_blank" className="fb"></a>
            </div>
            <div className="space-social"></div>
            <div className="social">
              <a href="https://www.twitter.com/testigo_social" target="_blank" className="tw"></a>
            </div>
          </div>
          <div className="col-md-6 mobile-logo">
            <div className="logo">
              <p>Una iniciativa de:</p>
              <a href="http://www.tm.org.mx" target="_blank" className="tm"></a>
            </div>
            <span data-toggle="modal" data-target="#privacyNotice">Aviso de Privacidad</span>
            <div className="logo">
              <p>Con el apoyo de:</p>
              <a href="https://www.gov.uk/government/world/organisations/british-embassy-mexico-city.es-419" target="_blank" className="uk"></a>
            </div>
          </div>
        </div>
      </div>
    );
  }
}

export default Footer;
