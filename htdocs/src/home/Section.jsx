import React from 'react';
import Chart from 'chart.js'
import {Link} from 'react-router';
import { getParameter, formatDate } from '../helpers.js';

class Home extends React.Component {
  constructor(props) {
    super(props);
  }

  componentDidMount() {
    let url = '/stats/gacm';
    let data = {};
    let charts = {};
    let ui = {
      totalContracts: $('span.totalContracts'),
      totalBudget: $('span#totalBudget'),
      totalAward: $('span#totalAward'),
      totalAmount: $('span#totalAmount'),
      firstDate: $('span#firstDate'),
      lastDate: $('span#lastDate'),
      description: $('span#orgDescription')
    }

    // Random hero photo
    let img = "url('images/hero_photo_"+ (Math.floor(Math.random() * 4) + 1) +".png')";
    $('#hero > div.photo').css( 'background-image', img );

    // Slider
    $('.carousel').carousel({
      interval: 8000,
      keyboard: false,
      pause: "hover",
    });

    // Dynamically set bucket used, default to 'gacm'
    if( getParameter('bucket') ) {
      url = '/stats/' + getParameter('bucket');
    }

    // Load stats
    $.ajax({
      type: "GET",
      url: url,
      success: function( res ) {
        data = JSON.parse(res);
      }
    }).done(function() {
      // Adjust labels
      ui.firstDate.text( formatDate( data.firstDate, 'LL' ) );
      ui.lastDate.text( formatDate( data.lastDate, 'LL' ) );
      ui.description.text( data.description );
      ui.totalContracts.text( data.contracts.total );
      ui.totalBudget.text((data.contracts.budget / 1000000).toLocaleString(undefined, {
        minimumFractionDigits: 1,
        maximumFractionDigits: 1
      }));
      ui.totalAward.text((data.contracts.awarded / 1000000).toLocaleString(undefined, {
        minimumFractionDigits: 1,
        maximumFractionDigits: 1
      }));
      ui.totalAmount.text( data.contracts.budget.toLocaleString(undefined, {
        minimumFractionDigits: 1,
        maximumFractionDigits: 1
      }));

      // Prepare chart data
      let directP = ((data.method['limited'].budget * 100) / data.contracts.budget).toFixed(2);
      let limitedP = ((data.method['selective'].budget * 100) / data.contracts.budget).toFixed(2);
      let publicP = ((data.method['open'].budget * 100) / data.contracts.budget).toFixed(2);
      let charts = {
        limited: {
          c: false,
          data: {
            labels: [
              'Adjudicación Directa (%)',
              'Total Contratado (%)'
            ],
            datasets: [
              {
                data: [directP, (100 - directP).toFixed(2)],
                backgroundColor: ["#CCB3FF", "#EEEEEE"],
              }
            ]
          }
        },
        selective: {
          c: false,
          data: {
            labels: [
              'Adjudicación Directa (%)',
              'Invitación a cuando menos 3 personas (%)',
              'Total Contratado (%)'
            ],
            datasets: [
              {
                data: [directP, limitedP, (100 - (Number(directP) + Number(limitedP))).toFixed(2)],
                backgroundColor: ["#CCB3FF", "#FF6384", "#EEEEEE"],
              }
            ]
          }
        },
        open: {
          c: false,
          data: {
            labels: [
              'Adjudicación Directa (%)',
              'Invitación a cuando menos 3 personas (%)',
              'Licitación Pública (%)'
            ],
            datasets: [
              {
                data: [directP, limitedP, publicP],
                backgroundColor: ["#CCB3FF", "#FF6384", "#7DE7CF"],
              }
            ]
          }
        }
      }

      // Configure data slider
      let slides = $('#content-slides');
      slides.on('slid.bs.carousel', function() {
        let active = slides.find('div.active');
        let method = active.data('section');
        if( method ) {
          active.find('span.contracts').text( data.method[method].total );
          active.find('span.amount').text( data.method[method].budget.toLocaleString({
            useGrouping: true
          }));

          if( !charts[method].c ) {
            charts[method].c = new Chart( active.find('canvas'), {
              type: "pie",
              data: charts[method].data,
              options: {
                responsive: true,
                responsiveAnimationDuration: 500,
                padding: 10
              }
            });
          }
        }
      });
    });
  }

  render() {
    return(
      <div>
        {/* Hero photo */}
        <div id="hero" className="row">
          <div className="photo col-md-12">
            <div className="logo"></div>
            <h2>El dinero público también es tu dinero</h2>
          </div>
        </div>

        {/* Facts */}
        <div id="facts" className="row inner-row">
          <div className="col-md-4">
            <p className="txt-centered">Número de <span className="txt-bold">procedimientos de contratación</span> registrados</p>
            <p className="highlight txt-centered txt-mono">
              <span className="counter totalContracts">0</span>
            </p>
          </div>
          <div className="col-md-4">
            <p className="txt-centered"><span className="txt-bold">Presupuesto asignado</span> de las contrataciones registradas</p>
            <p className="highlight txt-centered txt-mono">
              $<span id="totalBudget" className="counter">0</span>M
            </p>
          </div>
          <div className="col-md-4">
            <p className="txt-centered">Monto total <span className="txt-bold">contratado</span>a través de los procedimientos registrados</p>
            <p className="highlight txt-centered txt-mono">
              $<span id="totalAward" className="counter">0</span>M
            </p>
          </div>
        </div>

        {/* Highlights */}
        <div id="highlights" className="row inner-row">
          <div id="content-slides" className="col-md-12 carousel slide" data-ride="carousel">
            <ol className="carousel-indicators">
              <li data-target="#content-slides" data-slide-to="0" className="active"></li>
              <li data-target="#content-slides" data-slide-to="1"></li>
              <li data-target="#content-slides" data-slide-to="2"></li>
              <li data-target="#content-slides" data-slide-to="3"></li>
            </ol>

            <div className="carousel-inner" role="listbox">
              <div className="item active">
                <p className="info">
                  <span className="bg-green">Entre el <span id="firstDate"></span> y el <span id="lastDate"></span>, se han adjudicado</span>
                </p>
                <h1 className="txt-mono">$<span id="totalAmount" className="counter">0</span> MXN</h1>
                <p className="info">
                  <span className="bg-green">para <strong><span id="orgDescription"></span></strong>, mediante <span className="counter totalContracts">0</span> contratos.</span>
                </p>
              </div>
              <div className="item" data-section="limited">
                <div className="chart">
                  <canvas className="dataChart" width="460" height="320"></canvas>
                </div>
                <div className="info">
                  <h2>Adjudicación Directa</h2>
                  <h3>$ <span className="amount">0</span> MXN</h3>
                  <p>Han sido adjudicados directamente mediante <span className="contracts">0</span> contratos.</p>
                  <p className="details">En una adjudicación directa se entrega un contrato a una persona o empresa sin realizar un concurso público y abierto.</p>
                </div>
              </div>
              <div className="item" data-section="selective">
                <div className="chart">
                  <canvas className="dataChart" width="460" height="320"></canvas>
                </div>
                <div className="info">
                  <h2>Invitación a cuado menos 3 personas</h2>
                  <h3>$ <span className="amount">0</span> MXN</h3>
                  <p>Han sido adjudicados por invitación a cuando menos tres personas, mediante <span className="contracts">0</span> contratos.</p>
                  <p className="details">En una invitación a cuando menos tres personas se entrega un contrato mediante un concurso en el que solo participan un número restringido de personas o empresas, seleccionadas por la dependencia gubernamental que realiza la contratación.</p>
                </div>
              </div>
              <div className="item" data-section="open">
                <div className="chart">
                  <canvas className="dataChart" width="460" height="320"></canvas>
                </div>
                <div className="info">
                  <h2>Licitación Pública</h2>
                  <h3>$ <span className="amount">0</span> MXN</h3>
                  <p>Han sido adjudicados por licitación pública, mediante <span className="contracts">0</span> contratos.</p>
                  <p className="details">En una licitación pública se entrega un contrato mediante un concurso que está abierto a cualquier persona o empresa. En México existen licitaciones públicas nacionales, internacionales bajo tratados de libre comercio e internacionales abiertas.</p>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Main video */}
        <div id="video" className="row inner-row">
          <div className="col-md-12">
            <div className="holder"></div>
          </div>
          <div className="col-md-4 item">
            <Link to={'/contracts'}>
              <span className="btn-black txt-upper">Contratos</span>
              <p>Consulta la información alrededor de una <strong>contratación pública</strong>, desde su planeación hasta su ejecución y pago</p>
              <span className="icon contracts"></span>
            </Link>
          </div>
          <div className="col-md-4 item">
            <Link to={'/indicators'}>
              <span className="btn-black txt-upper">Indicadores</span>
              <p>Revisa <strong>estadísticas</strong> sobre el sistema de contrataciones de México y analiza cómo está funcionando</p>
              <span className="icon markers"></span>
            </Link>
          </div>
          <div className="col-md-4 item">
            <Link to={'/register'}>
              <span className="btn-black txt-upper">Notificaciones</span>
              <p>Recibe <strong>alertas e información</strong> oportuna sobre lo que ocurre en los procedimientos de contratación que te interesan</p>
              <span className="icon notifications"></span>
            </Link>
          </div>
          <div className="col-md-12">
            <br />
            <br />
            <span className="btn-black">Comienza a monitorear cómo se gasta tu dinero con TS 2.0</span>
          </div>
        </div>

        {/* Bottom text */}
        <div className="row info bg-blue">
          <div className="inner-row">
            <p>La apertura y la participación de la ciudadanía en las compras públicas se traduce en mejores bienes y servicios públicos para las <strong>comunidades</strong>, más oportunidades de negocio para emprendedores y <strong>empresas</strong> y una mayor rendición de cuentas de los <strong>gobiernos.</strong></p>
          </div>
        </div>
      </div>
    );
  }
}

export default Home;
