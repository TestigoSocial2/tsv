import React from 'react';
import {Link} from 'react-router';

class Menu extends React.Component {
  render() {
    return(
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
    );
  }
}

export default Menu;
