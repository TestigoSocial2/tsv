import React from 'react';
import countriesList from '../countries.js';

class Form extends React.Component {
  constructor(props) {
    super(props);
  }

  componentDidMount() {
    let agencies = $('div.agency-grid div.item');
    let selectedAgencies = $('div.agency-grid input#selectedAgencies');
    agencies.click(function() {
      // Toggle state
      let a = $(this);
      let lbl = a.find('span.label');
      a.toggleClass('selected');
      if( lbl.text().toLowerCase() == 'seguir' ) {
        lbl.text( 'siguiendo' );
      } else {
        lbl.text( 'seguir' );
      }

      // Adjust result
      let res = [];
      agencies.each( function( i, v ) {
        if( $(v).hasClass( 'selected') ) {
          res.push( $(v).data('value') );
        }
      });
      selectedAgencies.val(JSON.stringify(res));
    });

    let projects = $('div.project-grid div.item');
    let selectedProjects = $('div.project-grid input#selectedProjects');
    projects.click(function() {
      // Toggle state
      let p = $(this);
      let lbl = p.find('span.label');
      p.toggleClass('selected');
      if( lbl.text().toLowerCase() == 'seguir' ) {
        lbl.text( 'siguiendo' );
      } else {
        lbl.text( 'seguir' );
      }

      // Adjust result
      let res = [];
      projects.each( function( i, v ) {
        if( $(v).hasClass( 'selected') ) {
          res.push( $(v).data('value') );
        }
      });
      selectedProjects.val(JSON.stringify(res));
    });

    let form = $('form#registerForm');
    form.pixativeFormValidator({
      msgHolder: 'title',
      errorMessages: {
        required: 'El campo es requerido',
        minlength: 'El valor proporcionado debe ser de al menos %s caracteres',
        maxlength: 'El valor proporcionado debe ser de máximo %s caracteres',
        email: 'El valor porporcionado no parece ser una dirección de correo valida',
        integer: 'El valor proporcionado debe ser un número',
        phone: 'El valor porporcionado no parece ser un número telefonico valido'
      },
      onError: function() {
        $('.validator-error').tooltip('destroy');
        $('.validator-error').tooltip();
      },
      onSuccess: () => {
        // Prepare data
        let data = {}
        form.serializeArray().forEach(function(el) {
          if( el.value == 'on' ) {
            el.value = true;
          }
          if( el.value == 'off' ) {
            el.value = false;
          }
          data[el.name] = el.value;
        });
        data.selectedAgencies = JSON.parse(data.selectedAgencies);
        data.selectedProjects = JSON.parse(data.selectedProjects);
        this.props.onSubmit(data);
      }
    });
  }

  render() {
    let ages = [];
    for( let i = 15; i < 100; i++ ) {
      ages.push(<option value={i}>{i}</option>);
    }

    let countries = [];
    countriesList.forEach((c)=>{
      countries.push(<option value={c.code}>{c.name}</option>);
    });
    return(
      <div id="setup">
        <form id="registerForm">
          <div className="row">
            <div className="col-md-3">
              <h2 className="number-title"><span>1</span>Registro</h2>
              <p className="txt-italic">Completa la siguiente información</p>
            </div>
            <div className="col-md-9">
              <table className="table table-striped green">
                <tbody>
                  <tr>
                    <td className="txt-upper col-md-2"><label>Correo Electrónico</label></td>
                    <td className="col-md-4">
                      <input type="email"
                             className="form-control"
                             id="user"
                             name="user"
                             placeholder="Email"
                             data-validator-required="true"
                             data-validator-email="true" />
                    </td>
                    <td className="txt-italic">La cuenta de correo electrónico registrada será tu nombre de usuario</td>
                  </tr>
                  <tr>
                    <td className="txt-upper col-md-2"><label>Contraseña</label></td>
                    <td className="col-md-4">
                      <input type="password"
                             className="form-control"
                             id="password"
                             name="password"
                             placeholder="Contraseña"
                             data-validator-required="true"
                             data-validator-minlength="6" />
                    </td>
                    <td className="txt-italic">La contraseña deberá tener una extensión de al menos 6 caracteres y contener al menos un número</td>
                  </tr>
                  <tr>
                    <td className="txt-upper col-md-2"><label>Tipo de Usuario</label></td>
                    <td className="col-md-4">
                      <select id="userType" name="userType" className="form-control" data-validator-required="true">
                        <option value="-" disabled="disaled">Selecciona una de las siguientes opciones</option>
                        <option value="public">Ciudadano</option>
                        <option value="legal">Legislador</option>
                        <option value="gov">Funcionario Público</option>
                        <option value="media">Periodista</option>
                        <option value="startup">Emprendedor</option>
                        <option value="business">Empresario</option>
                        <option value="other">Otro</option>
                      </select>
                    </td>
                    <td className="txt-italic">Selecciona una de las siguientes opciones</td>
                  </tr>
                  <tr>
                    <td className="txt-upper col-md-2"><label>Edad</label></td>
                    <td className="col-md-4">
                      <select id="age" name="age" className="form-control" data-validator-required="true">
                        {ages}
                      </select>
                    </td>
                    <td className="txt-italic">Selecciona tu edad (solo aplica para periodista, emprendedor, legislador, funcionario público, ciudadano)</td>
                  </tr>
                  <tr>
                    <td className="txt-upper col-md-2"><label>Código Postal</label></td>
                    <td className="col-md-4">
                      <input type="text"
                             className="form-control"
                             id="postalCode"
                             name="postalCode"
                             data-validator-required="true"
                             data-validator-integer="true" />
                    </td>
                    <td className="txt-italic">Introduce tu código postal</td>
                  </tr>
                  <tr>
                    <td className="txt-upper col-md-2"><label>País</label></td>
                    <td className="col-md-4">
                      <select id="country" name="country" className="form-control" data-validator-required="true">
                        {countries}
                      </select>
                    </td>
                    <td className="txt-italic">Selecciona tu país de residencia</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
          <div className="row">
            <div className="col-md-3">
              <h2 className="number-title"><span>2</span>Proyectos</h2>
              <p className="txt-italic">Selecciona las agencias, proyectos o sectores que te interesan</p>
            </div>
            <div className="row">
              <div className="agency-grid col-md-4">
                <h4 className="txt-upper">Agencias</h4>
                <input type="text" id="selectedAgencies" name="selectedAgencies" value="[]" className="hidden" />
                <div className="item" data-value="cdmx">
                  <div className="logo">
                    <img src="https://twitter.com/gobcdmx/profile_image?size=original" />
                  </div>
                  <span className="label">Seguir</span>
                </div>
                <div className="item" data-value="a2">
                  <div className="logo">
                    <img src="https://twitter.com/nvoaeropuertomx/profile_image?size=original" />
                  </div>
                  <span className="label">Seguir</span>
                </div>
                <div className="clear"></div>
              </div>
              <div className="project-grid col-md-3">
                <h4 className="txt-upper">Proyectos</h4>
                <input type="text" id="selectedProjects" name="selectedProjects" value="[]" className="hidden" />
                <div className="item" data-value="p1">
                  <h5>NAICM</h5>
                  <div className="img">
                    <img src="https://twitter.com/nvoaeropuertomx/profile_image?size=original" />
                  </div>
                  <span className="label">Seguir</span>
                  <p>Construcción del nuevo aeropuerto internacional de la ciduad de México.</p>
                </div>
              </div>
              <div className="col-md-2">
                <h4 className="txt-upper">Sectores</h4>
                <ul>
                  <li>
                    <span className="glyphicon glyphicon-home" aria-hidden="true"></span>Infraestructura
                  </li>
                </ul>
              </div>
            </div>
          </div>
          <div className="row">
            <div className="col-md-3">
              <h2 className="number-title"><span>3</span>Contacto</h2>
              <p className="txt-italic">Selecciona los canales por los que Testigo Social 2.0 puede enviarte notificaciones e información.</p>
            </div>
            <div className="col-md-9">
              <table className="table table-striped blue">
                <tbody>
                  <tr>
                    <td className="txt-upper col-md-2"><label>Correo Electrónico</label></td>
                    <td className="col-md-4">
                      <input type="email"
                             className="form-control"
                             id="notificationEmail"
                             name="notificationEmail"
                             data-validator-email="true" />
                    </td>
                    <td>
                      <input id="enableEmailNotifications"
                             name="enableEmailNotifications"
                             type="checkbox"
                             checked="checked" />
                    </td>
                    <td className="txt-italic">Quiero recibir notificaciones a través de un correo electrónico</td>
                  </tr>
                  <tr>
                    <td className="txt-upper col-md-2"><label>SMS</label></td>
                    <td className="col-md-4">
                      <input type="text"
                             className="form-control"
                             id="notificationSMS"
                             name="notificationSMS"
                             data-validator-integer="true"
                             data-validator-minlength="10"
                             data-validator-maxlength="10" />
                    </td>
                    <td>
                      <input id="enableSMSNotifications"
                             name="enableSMSNotifications"
                             type="checkbox"
                             checked="checked" />
                    </td>
                    <td className="txt-italic">Quiero recibir notificaciones a través de mensajes a un número celular (sin costo)</td>
                  </tr>
                  <tr>
                    <td className="txt-upper col-md-2"><label>Facebook</label></td>
                    <td className="col-md-4">
                      <input disabled="disabled" type="text" className="form-control" id="notificationFB" />
                    </td>
                    <td>
                      <input disabled="disabled" id="enableFB" type="checkbox" />
                    </td>
                    <td className="txt-italic">Quiero recibir notificaciones e información a través mensajes en Facebook Messenger</td>
                  </tr>
                  <tr>
                    <td className="txt-upper col-md-2"><label>Twitter</label></td>
                    <td className="col-md-4">
                      <input disabled="disabled" type="text" className="form-control" id="notificationTW" />
                    </td>
                    <td>
                      <input disabled="disabled" id="enableTW" type="checkbox" />
                    </td>
                    <td className="txt-italic">Quiero recibir notificaciones a través de mensajes directos y notificaciones en Twitter</td>
                  </tr>
                </tbody>
              </table>
              <p>
                Consulta nuestro aviso de privacidad <a href="#" data-toggle="modal" data-target="#privacyNotice">aquí.</a>
              </p>
              <div className="checkbox">
                <label>
                  <input id="acceptPrivacyTerms" type="checkbox" checked="checked" /> Acepto las politicas de privacidad
                </label>
              </div>
              <button type="submit" className="btn btn-black btn-lg">Comienza a recibir notificaciones</button>
            </div>
          </div>
        </form>
      </div>
    );
  }
}

export default Form;
