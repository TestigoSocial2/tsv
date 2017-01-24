import React from 'react';
import { getParameter, formatAmount, formatDate } from './helpers.js'

class Description extends React.Component {
  render() {
    return(
      <div className="row info bg-blue">
        <div className="inner-row">
          <h2>{this.props.title}</h2>
          <p>{this.props.content}</p>
        </div>
      </div>
    );
  }
}

// <div class="item-cell">
//   <p class="lbl"><span class="index">1</span> 13099024</p>
//   <p><a href="#">Nivelación y limpieza del terreno del Nuevo Aeropuerto Internacional de la Ciudad de México.</a></p>
//   <div>
//     <div class="left">
//       <p class="lbl">Monto</p>
//       <p><strong>$1,763,500,240.00 MXN</strong></p>
//       <p class="lbl">Proveedor</p>
//       <p>COCONAL SAPI DE CV</p>
//     </div>
//     <div class="right">
//       <p class="lbl">Comprador</p>
//       <p>GACM</p>
//       <span class="btn-black active">Ver contrato</span>
//     </div>
//     <div class="clear"></div>
//   </div>
// </div>
class TableItem extends React.Component {
  constructor(props) {
    super(props);
    this.handleClick = this.handleClick.bind(this);
  }

  handleClick(e) {
    e.preventDefault();
    this.props.onClick( this.props.contract );
  }

  render() {
    return (
      <tr>
        <td width="50%">
          <p className="lbl">{this.props.contract.releases[0].tender.id}</p>
          <a onClick={this.handleClick}>{this.props.contract.releases[0].tender.description}</a>
        </td>
        <td>
          <p className="lbl">Monto</p>
          <h4>{formatAmount( this.props.contract.releases[0].tender.value.amount || 0 )}</h4>
          <p className="lbl">Proveedor</p>
          <p>{this.props.contract.releases[0].awards[0].suppliers[0].name}</p>
        </td>
        <td>
          <p className="lbl">Comprador</p>
          <p>{this.props.contract.releases[0].buyer.name}</p>
        </td>
        <td width="15%">
          <a className="btn-black active" onClick={this.handleClick}>Ver Contrato</a>
        </td>
      </tr>
    );
  }
}

class SearchBar extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      filter: null,
      value: null,
      limit: 30
    }
  }

  componentDidMount() {
    this.ui = {}
    this.ui.form = $('form#queryForm');
    this.ui.bullets = this.ui.form.find('div.bullets span');
    this.ui.input = this.ui.form.find( 'input#query' );

    // Setup calendar
    this.ui.form.find("span[data-filter='date']").datepicker({
      clearBtn: true,
      assumeNearbyYear: true,
      format: 'mm/dd/yyyy',
      language: 'es',
      maxViewMode: 2,
      multidate: 2,
      todayHighlight: true
    }).on( 'hide', (function( e ) {
      e.dates.sort(function(a, b) {
        return new Date(a).getTime() - new Date(b).getTime();
      });
      let lbl = moment(e.dates[0]).format('MMMM Do YYYY');
      let val = moment(e.dates[0]).format('MM-DD-YYYY');
      if( e.dates.length > 1 ) {
        lbl += ' a ' + moment(e.dates[1]).format('MMMM Do YYYY')
        val += '|' + moment(e.dates[1]).format('MM-DD-YYYY');
      }

      // Update state
      this.setState({
        value: val
      });

      // Update UI
      this.ui.input.val( lbl );
      this.ui.input.focus();
    }).bind(this));

    // Setup slider
    this.ui.form.find("span[data-filter='amount']").popover({
      html: true,
      title: 'Seleccione el rango a utilizar como filtro (MXN)',
      content: '<b>$0</b><input id="amountSlider" type="text" /><b>$100,000,000</b>',
      placement: 'bottom',
      trigger: 'focus'
    }).on( 'shown.bs.popover', (function() {
      $("#amountSlider").slider({
        step: 50000,
        min: 0,
        max: 100000000,
        value: [20000000,80000000],
        formatter: function( value ) {
          if( Array.isArray( value ) ) {
            var lbl = '$' + value[0].toLocaleString() + ' a ' + '$' + value[1].toLocaleString();
            return lbl;
          }
          return '';
        }
      }).on( 'slide', (function( e ) {
        this.setState({
          value: e.value.join('|')
        });
        this.ui.input.val( '$' + e.value[0].toLocaleString() + ' a ' + '$' + e.value[1].toLocaleString() );
      }).bind(this));
    }).bind(this));

    // Update state when filter value changes
    this.ui.input.on( 'keyup', (function() {
      this.setState({
        value: this.ui.input.val()
      })
    }).bind(this));

    // Handle filter selection
    this.ui.bullets.on( 'click', (function( e ) {
      // Toggle active
      let target = $( e.target );
      this.ui.bullets.removeClass( 'active' );
      target.addClass( 'active' );

      // Update state
      this.setState({
        filter: target.data('filter')
      });

      switch (target.data( 'filter' )) {
        case 'date':
          target.datepicker( 'show' );
          break;
        case 'amount':
          target.popover( 'toggle' );
          break;
      }
    }).bind(this) );

    // Handle form submit
    this.ui.form.on( 'submit', (function( e ) {
      e.preventDefault();
      this.props.onSubmit(this.state);
    }).bind(this));
  }

  render() {
    return (
      <div className="inner-row">
        <div className="row">
          <div className="col-md-12">
            <h2>Buscador de Contratos</h2>
            <p>Encuentra contratos o procedimientos de contratación registrados en Testigo Social 2.0 haciendo uso de los distintos filtros de búsqueda disponibles.</p>
            <form id="queryForm">
              <div className="input-group">
                <input type="text" className="form-control" id="query" name="query" placeholder="Buscar..." />
                <span className="input-group-btn">
                  <button className="btn btn-primary" type="submit">Buscar</button>
                </span>
              </div>
              <div className="bullets">
                <span className="btn-black active" data-filter="date">Fecha</span>
                <span className="btn-black" data-filter="amount">Monto</span>
                <span className="btn-black" data-filter="buyer">Comprador</span>
                <span className="btn-black" data-filter="provider">Proveedor</span>
                <span className="btn-black" data-filter="procedureType">Tipo de Procedimiento</span>
                <span className="btn-black" data-filter="procedureNumber">No. de Procedimiento</span>
                <span className="btn-black" data-filter="contractNumber">No. de Contrato</span>
              </div>
            </form>
          </div>
        </div>
      </div>
    );
  }
}

class SearchResults extends React.Component {
  constructor(props) {
    super(props);
    this.itemSelected = this.itemSelected.bind(this);
  }

  itemSelected(contract) {
    this.props.onSelection(contract);
  }

  render() {
    let content = (<tr><td><h4><i>Sin resultados que mostrar...</i></h4></td></tr>);
    if( this.props.items.length > 0 ) {
      content = this.props.items.map((contract) =>
        <TableItem onClick={this.itemSelected} key={contract.releases[0].id} contract={contract} />
      );
    }

    return(
      <div className="inner-row">
        <div className="row">
          <div className="col-md-12">
            <h2>Resultados de la busqueda</h2>
            <p>Estos son los procedimientos de contratación registrados en Testigo Social 2.0 que satisfacen los creiterios de selección proporcionados</p>
            <table className="table table-striped latest-contracts">
              <tbody>
                {content}
              </tbody>
            </table>
          </div>
        </div>
      </div>
    );
  }
}

class Details extends React.Component {
  constructor(props) {
    super(props);
    this.onClose = this.onClose.bind(this);
  }

  componentDidMount() {
    let bullets = $('div.contract-header div.bullets a');
    bullets.on( 'click', function( e ) {
      bullets.removeClass('active');
      $(e.target).addClass('active');
    });
  }

  onClose(e) {
    e.preventDefault();
    this.props.onClose();
  }

  render() {
    let release = this.props.contract.releases[0];
    return (
      <div className="inner-row contract-details">
        {/* Header */}
        <div className="row contract-header">
          <div className="col-md-12">
            <div>
              <div className="bg-gray">
                <h2 className="block-title">{release.tender.title}</h2>
                <h3>{formatAmount( release.tender.value.amount || 0 )}</h3>
                <p>{release.tender.description}</p>
              </div>
              <div className="bullets">
                <a href="#planning" aria-controls="planning" data-toggle="tab" className="btn-black active">Planeación</a>
                <a href="#tender" aria-controls="tender" data-toggle="tab" className="btn-black">Licitación</a>
                <a href="#award" data-toggle="tab" className="btn-black">Adjudicación</a>
                <a href="#contract" data-toggle="tab" className="btn-black">Contratación</a>
                <a href="#implementation" data-toggle="tab" className="btn-black disabled" disabled="disabled">Implementación</a>
                <a onClick={this.onClose}>Volver al listado de Resultados</a>
              </div>
            </div>
          </div>
        </div>

        {/* Content */}
        <div className="tab-content">
          {/* Planning */}
          <div role="tabpanel" className="tab-pane active fade in" id="planning">
            <div className="row contract-highlights border-bottom">
              <div className="col-md-4">
                <p className="lbl">Fuente Presupuestaria</p>
                <h4>{release.planning.budget.source}</h4>
              </div>
              <div className="col-md-4">
                <p className="lbl">Identificador del Presupuesto</p>
                <h4>{release.planning.budget.id}</h4>
              </div>
              <div className="col-md-4">
                <p className="lbl">Monto Asignado</p>
                <h4>{formatAmount(release.planning.budget.amount.amount || 0)}</h4>
              </div>
            </div>

            <div className="row border-bottom">
              <div className="col-md-12">
                <p className="lbl">Fundamento</p>
                <p className="txt-bold">{release.planning.budget.description}</p>
              </div>
            </div>

            <div className="row contract-highlights border-bottom">
              <div className="col-md-4">
                <p className="lbl">¿Proyecto Plurianual?</p>
                <h4>{'?'}</h4>
              </div>
              <div className="col-md-4">
                <p className="lbl">Proyecto Presupuestario</p>
                <h4>{release.planning.budget.project}</h4>
              </div>
              <div className="col-md-4">
                <p className="lbl">Identificador</p>
                <h4>{this.props.contract.releases[0].planning.budget.projectID}</h4>
              </div>
            </div>

            <div className="row border-bottom">
              <div className="col-md-12">
                <p className="lbl">Enlace a la información presupuestaria</p>
                <a href={this.props.contract.releases[0].planning.budget.uri}>
                  {this.props.contract.releases[0].planning.budget.uri}
                </a>
              </div>
            </div>
          </div>

          {/* Tender */}
          <div role="tabpanel" className="tab-pane fade" id="tender">
            <div className="row contract-highlights border-bottom">
              <div className="col-md-8">
                <p className="lbl">Título de la Licitación</p>
                <h4>{release.tender.title}</h4>
              </div>
              <div className="col-md-4">
                <p className="lbl">Identificador de la Licitación</p>
                <h4>{release.tender.id}</h4>
              </div>
            </div>

            <div className="row contract-highlights border-bottom">
              <div className="col-md-4">
                <p className="lbl">Tipo de Contratación</p>
                <h4>{release.tender.metodoDeAdquisicion}</h4>
              </div>
              <div className="col-md-8">
                <p className="lbl">Descripción de la Licitación</p>
                <h4>{release.tender.description}</h4>
              </div>
            </div>

            <div className="row contract-highlights border-bottom">
              <div className="col-md-4">
                <p className="lbl">Estatus de la Licitación</p>
                <h4>{release.tender.status}</h4>
              </div>
              <div className="col-md-4">
                <p className="lbl">Valor Máximo</p>
                <h4>{formatAmount(release.tender.value.amount || 0)}</h4>
              </div>
              <div className="col-md-4">
                <p className="lbl">Método por el que se realiza</p>
                <h4>{release.tender.procurementMethod}</h4>
              </div>
            </div>

            <div className="row contract-highlights border-bottom">
              <div className="col-md-4">
                <p className="lbl">Cáracter del Proceso</p>
                <h4>
                  {release.tender.submissionMethod ? release.tender.submissionMethod.join(',') : '-'}
                </h4>
              </div>
              <div className="col-md-4">
                <p className="lbl">Forma del Proceso</p>
                <h4>{'?'}</h4>
              </div>
              <div className="col-md-4">
                <p className="lbl">Criterio de Adjudicación</p>
                <h4>{release.tender.procurementMethodRationale}</h4>
              </div>
            </div>

            <div className="row contract-highlights border-bottom">
              <div className="col-md-4">
                <p className="lbl">Periodo de Recepción de Propuestas</p>
                <h4>{'?'}</h4>
              </div>
              <div className="col-md-4">
                <p className="lbl">Periodo de Aclaraciones</p>
                <h4>{'?'}</h4>
              </div>
              <div className="col-md-4">
                <p className="lbl">Aclaraciones</p>
                <h4>{'?'}</h4>
              </div>
            </div>

            <div className="row contract-highlights border-bottom">
              <div className="col-md-4">
                <p className="lbl">Testigo Social</p>
                <h4>{'?'}</h4>
              </div>
              <div className="col-md-8">
                <p className="lbl">Criterio de Elegibilidad</p>
                <p className="txt-bold">{release.tender.awardCriteria}</p>
              </div>
            </div>

            <div className="row contract-highlights border-bottom">
              <div className="col-md-4">
                <p className="lbl">No. de Propuestas Recibidas</p>
                <h4>{release.tender.numberOfTenderers}</h4>
              </div>
              <div className="col-md-4">
                <p className="lbl">Fecha de Adjudicación</p>
                <h4>{formatDate(release.date, 'DD/MM/YYYY')}</h4>
              </div>
              <div className="col-md-4">
                <p className="lbl">Identificador de Entidad</p>
                <h4>{this.props.contract.publisher.uid}</h4>
              </div>
            </div>
          </div>

          {/* Award */}
          <div role="tabpanel" className="tab-pane fade" id="award">
            <div className="row contract-highlights border-bottom">
              <div className="col-md-8">
                <p className="lbl">Título de la Adjudicación</p>
                <h4>{'?'}</h4>
              </div>
              <div className="col-md-4">
                <p className="lbl">Identificador de la Adjudicación</p>
                <h4>{release.awards[0].id}</h4>
              </div>
            </div>

            <div className="row contract-highlights border-bottom">
              <div className="col-md-4">
                <p className="lbl">Estatus</p>
                <h4>{release.awards[0].status}</h4>
              </div>
              <div className="col-md-8">
                <p className="lbl">Descripción de la Adjudicación</p>
                <p className="txt-bold">{'?'}</p>
              </div>
            </div>

            <div className="row contract-highlights border-bottom">
              <div className="col-md-4">
                <p className="lbl">Fecha</p>
                <h4>{'?'}</h4>
              </div>
              <div className="col-md-4">
                <p className="lbl">Valor y Moneda</p>
                <h4>{formatAmount(release.awards[0].value.amount)}</h4>
              </div>
              <div className="col-md-4">
                <p className="lbl">No. de Inconformidades Recibidas</p>
                <h4>{'?'}</h4>
              </div>
            </div>

            <div className="row contract-highlights border-bottom">
              <div className="col-md-4">
                <p className="lbl">No. de Inconformidades Rechazadas</p>
                <h4>{'?'}</h4>
              </div>
              <div className="col-md-4">
                <p className="lbl">Identificador de Proveedor</p>
                <h4>{release.awards[0].suppliers[0].identifier.id}</h4>
              </div>
              <div className="col-md-4">
                <p className="lbl">Nombre de Proveedor</p>
                <h4>{release.awards[0].suppliers[0].identifier.legalName}</h4>
              </div>
            </div>
          </div>

          {/* Contracting */}
          <div role="tabpanel" className="tab-pane fade" id="contract">
            <div className="row contract-highlights border-bottom">
              <div className="col-md-4">
                <p className="lbl">Identificador del Contrato</p>
                <h4>{release.contracts[0].id}</h4>
              </div>
              <div className="col-md-8">
                <p className="lbl">Título del Contrato</p>
                <h4>{release.contracts[0].title}</h4>
              </div>
            </div>

            <div className="row contract-highlights border-bottom">
              <div className="col-md-8">
                <p className="lbl">Descripción del Contrato</p>
                <p className="txt-bold">{'?'}</p>
              </div>
              <div className="col-md-4">
                <p className="lbl">Estatus</p>
                <h4>{release.contracts[0].status}</h4>
              </div>
            </div>

            <div className="row contract-highlights border-bottom">
              <div className="col-md-4">
                <p className="lbl">Periodo</p>
                <h4>
                  {
                    formatDate(release.contracts[0].period.startDate, 'DD/MM/YYYY')
                    + ' - ' +
                    formatDate(release.contracts[0].period.endDate, 'DD/MM/YYYY')
                  }
                </h4>
              </div>
              <div className="col-md-4">
                <p className="lbl">Valor y Moneda</p>
                <h4>{formatAmount(release.contracts[0].value.amount)}</h4>
              </div>
              <div className="col-md-4">
                <p className="lbl">Fecha de Firma del Contrato</p>
                <h4>{formatDate(release.contracts[0].dateSigned, 'DD/MM/YYYY')}</h4>
              </div>
            </div>

            <div className="row contract-highlights border-bottom">
              <div className="col-md-12">
                <p className="lbl">¿SE MODIFICÓ EL CONTRATO?</p>
                <h4>{'?'}</h4>
              </div>
            </div>
          </div>

          {/* Implementation */}
          <div role="tabpanel" className="tab-pane fade" id="implementation">
            <h4><i>Sin información que mostrar por el momento...</i></h4>
          </div>
        </div>
      </div>
    );
  }
}

class Section extends React.Component {
  constructor(props) {
    super(props);
    this.runQuery = this.runQuery.bind(this);
    this.showDetails = this.showDetails.bind(this);
    this.hideDetails = this.hideDetails.bind(this);
    this.state = {
      items: [],
      details: null
    };
  }

  runQuery(query) {
    let url = '/query/gacm';
    if( getParameter('bucket') ) {
      url = '/query/' + getParameter('bucket');
    }

    // Submit query and update component state with results
    $.ajax({
      type: "POST",
      url: url,
      data: {
        query: JSON.stringify(query)
      },
      success: (function( res ) {
        this.setState({
          items: JSON.parse( res )
        });
      }).bind(this)
    });
  }

  showDetails(contract) {
    this.setState({
      details: contract
    });
  }

  hideDetails() {
    this.setState({
      details: null
    });
  }

  render() {
    if( !this.state.details ) {
      return (
        <div>
          <Description
            title="Contratos"
            content="Consulta cada contrato que está registrado en Testigo Social 2.0. Podrás encontrar información para cada una de las etapas del procedimiento de contratación, desde su planeación hasta su implementación." />
          <SearchBar onSubmit={this.runQuery} />
          <SearchResults items={this.state.items} onSelection={this.showDetails} />
        </div>
      );
    } else {
      return (
        <div>
          <Description
            title="Contratos"
            content="En esta sección podrás visualizar la información completa y agregada de un contrato, desde su planeación hasta su implementación. Como base se utiliza el Estándar de Datos para las Contrataciones Abiertas." />
          <Details contract={this.state.details} onClose={this.hideDetails} />
        </div>
      );
    }
  }
}

export default Section;
