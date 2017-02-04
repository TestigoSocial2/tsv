import React from 'react';
import {Link} from 'react-router';

class Main extends React.Component {
  componentDidMount() {
    $('[data-toggle="tooltip"]').tooltip();
  }

  render() {
    return(
      <div className="container-fluid">
        {/* Top Menu */}
        <div id="top" className="row">
          <div className="col-md-12">
            <a href="index.html" className="logo"></a>

            <ul>
              <div className="clear"></div>
              <li>
                <Link to={'/'}>¿Qué es TS 2.0?</Link>
              </li>
              <li>
                <Link to={'/contracts'}>Contratos</Link>
              </li>
              <li>
                <Link to={'/indicators'}>Indicadores</Link>
              </li>
              <li>
                <Link to={'/register'}>Notificaciones</Link>
              </li>
              <li><a href="http://www.tm.org.mx" target="_blank">Contrataciones Abiertas</a></li>
            </ul>

            <div className="social">
              <a href="https://www.twitter.com/testigo_social" className="tw"></a>
              <a href="https://www.facebook.com/testigosocial" className="fb"></a>
            </div>
          </div>
        </div>

        {/* Routes */}
        <div id="route">
          {this.props.children}
        </div>

        {/* Footer */}
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
      </div>
    );
  }
}

export default Main;
