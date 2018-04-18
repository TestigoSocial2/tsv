import React from 'react';
import {Link} from 'react-router';
var FontAwesome = require('react-fontawesome');

class Menu extends React.Component {

  componentDidMount() {
    $('.navbar-collapse').collapse('hide')
    let links = $('#top ul a');
    links.click((e) => {
      links.removeClass('active');
      $(e.target).addClass('active');
    });
  }

  render() {
    return(
      <div id="top" className="row">
        <div className="col-md-12">
          <Link to={'/'} className="logo"></Link>

          <ul>
            <div className="clear"></div>
            <li>
              <Link to={'/informacion'}>¿Qué es TS 2.0?</Link>
            </li>
            <li>
              <Link to={'/contratos'}>Contratos</Link>
            </li>
            <li>
              <Link to={'/indicadores'}>Indicadores</Link>
            </li>
            <li>
              <Link to={'/registro'}>Notificaciones</Link>
            </li>
            <li><a href="http://www.contratacionesabiertas.mx" target="_blank">Contrataciones Abiertas</a></li>
          </ul>
          <div className="social">
            <a href="https://www.twitter.com/testigo_social" className="tw"></a>
            <a href="https://www.facebook.com/testigosocial" className="fb"></a>
          </div>
        </div>
        <nav className="navbar navbar-default">
            <div className="container-fluid">
              <div className="navbar-header">
                <Link to={'/'} className="logo2"></Link>
                <button style={{ margin : 0,padding : 0,border : 0}} className="navbar-toggle collapsed" data-toggle="collapse" data-target=".navbar-collapse">
                  <FontAwesome name='bars' aria-expanded="false"/>
                </button>
              </div>
              <div className="collapse navbar-collapse" id="menu">
                <ul className="nav navbar-nav">
                  <li>
                    <Link to={'/informacion'} data-toggle="collapse" data-target=".navbar-collapse">¿Qué es TS 2.0?</Link>
                  </li>
                  <li>
                    <Link to={'/contratos'} data-toggle="collapse" data-target=".navbar-collapse">Contratos</Link>
                  </li>
                  <li>
                    <Link to={'/indicadores'} data-toggle="collapse" data-target=".navbar-collapse">Indicadores</Link>
                  </li>
                  <li>
                    <Link to={'/registro'} data-toggle="collapse" data-target=".navbar-collapse">Notificaciones</Link>
                  </li>
                  <li><a href="http://www.contratacionesabiertas.mx/" target="_blank" data-toggle="collapse" data-target=".navbar-collapse">Contrataciones Abiertas</a>
                  </li>
                </ul>
              </div>
            </div>
          </nav>
      </div>
    );
  }
}

export default Menu;
