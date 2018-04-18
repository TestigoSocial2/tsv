import React from 'react';
import { formatDate } from '../helpers.js';

class SearchBar extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      query: {},
      limit: 30
    }
  }

  componentDidMount() {
    this.ui = {};
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
    }).on({
      'show' : (e) => {
        this.ui.form.find("span[data-filter='amount']").popover('hide');
      },
      'hide': (e) => {
        e.dates.sort(function(a, b) {
          return new Date(a).getTime() - new Date(b).getTime();
        });
        let lbl = formatDate(e.dates[0], 'MMMM Do YYYY');
        let q = {
          'releases.date':{
            '$gte': e.dates[0]
          }
        };
        if( e.dates.length > 1 ) {
          lbl += ' a ' + formatDate(e.dates[1], 'MMMM Do YYYY');
          q['releases.date']['$lte'] = e.dates[1];
        }

        // Update state
        this.setState({query: q});

        // Update UI
        this.ui.input.val( lbl );
        this.ui.input.focus();
      },
      'changeDate': () => {
        // Visually mark the selected range
        if( $('div.datepicker-days tbody tr td.active').length == 2 ) {
          let mark = false;
          let days = $('div.datepicker-days tbody td.active');
          $('div.datepicker-days tbody td').each(function(i, d) {
            if( d == days[0] ) {
              mark = true;
            }
            if( mark ) {
              $(d).addClass('range');
            }
            if( d == days[1] ) {
              mark = false;
            }
          });
        }
      }
    });

    // Setup slider
    this.ui.form.find("span[data-filter='amount']").popover({
      html: true,
      title: 'Seleccione el rango a utilizar como filtro (MXN)',
      content: '<b>$0 - $100,000,000</b><input id="amountSlider" type="text" />',
      placement: 'bottom',
      trigger: 'focus'
    }).on( 'shown.bs.popover', () => {
      this.ui.form.find("span[data-filter='date']").datepicker('hide');
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
      }).on( 'slide', (e) => {
        this.setState({
          query: {
            'releases.planning.budget.amount.amount': {
              '$gte': e.value[0],
              '$lte': e.value[1]
            }
          }
        });
        this.ui.input.val( '$' + e.value[0].toLocaleString() + ' a ' + '$' + e.value[1].toLocaleString() );
      });
    });

    // Handle filter selection
    this.ui.bullets.on( 'click', (function( e ) {
      // Toggle active
      let target = $( e.target );
      this.ui.bullets.removeClass( 'active' );
      target.addClass( 'active' );

      switch (target.data( 'filter' )) {
        case 'date':
          target.datepicker( 'show' );
          break;
        case 'amount':
          target.popover( 'toggle' );
          break;
      }
    }).bind(this));

    // Handle form submit
    this.ui.form.on( 'submit', (function( e ) {
      e.preventDefault();
      let q = this.state.query;
      let target = $( this.ui.form.find('div.bullets span.active') );
      switch (target.data( 'filter' )) {
        case 'provider':
          q = {
            'releases.awards.suppliers.identifier.legalName': {
              '$regex': '.*' + this.ui.input.val() + '.*'
            }
          };
          break;
        case 'buyer':
          q = {
            'releases.buyer.identifier.legalName': {
              '$regex': '.*' + this.ui.input.val() + '.*'
            }
          };
          break;
        case 'procedureType':
          break;
        case 'procedureNumber':
          q = { 'releases.ocid': this.ui.input.val() };
          break;
        case 'contractNumber':
          q = { 'releases.contracts.id': this.ui.input.val() };
          break;
      }
      this.setState({query: q});
      this.props.onSubmit(this.state);
    }).bind(this));

    // Clear filters
    this.ui.form.on( 'reset', (e) => {
      e.preventDefault();
      this.setState({
        filter: null,
        value: null,
        limit: 30
      });
      this.ui.input.val( '' );
    });
  }

  render() {
    return (
      <div className="inner-row">
        <div className="row">
          <div className="col-md-12">
            <h2 className="mobile-search">Buscador de Contratos</h2>
            <p className="mobile-search">Encuentra contratos o procedimientos de contratación registrados en Testigo Social 2.0 haciendo uso de los distintos filtros de búsqueda disponibles.</p>
            <form id="queryForm">
              <div className="input-group hide-contracts">
                <input type="text" className="form-control" id="query" name="query" placeholder="Buscar..." />
                <span className="input-group-btn">
                  <button className="btn" type="reset">Limpiar</button>
                  <button className="btn btn-primary" type="submit">Buscar</button>
                </span>
              </div>
              <div className="bullets hidden-sm hidden-xs" >
                <span className="btn-black active" data-filter="date">Fecha</span>
                <span className="btn-black" data-filter="amount">Monto</span>
                <span className="btn-black" data-filter="buyer">Comprador</span>
                <span className="btn-black" data-filter="provider">Proveedor</span>
                <span className="btn-black" data-filter="procedureType">Tipo de Procedimiento</span>
                <span className="btn-black" data-filter="procedureNumber">No. de Procedimiento</span>
                <span className="btn-black" data-filter="contractNumber">No. de Contrato</span>
              </div>
              <div className="bullets visible-sm visible-xs">
                <div className="col-sm-6 col-xs-6">
                  <span className="btn-black active" data-filter="date">Fecha</span>
                  <span className="btn-black" data-filter="buyer">Comprador</span>
                  <span className="btn-black" data-filter="procedureType">Tipo de Procedimiento</span>
                </div>
                <div className="col-sm-6 col-xs-6">
                  <span className="btn-black" data-filter="amount">Monto</span>
                  <span className="btn-black" data-filter="provider">Proveedor</span>
                  <span className="btn-black" data-filter="procedureNumber">No. de Procedimiento</span>
                </div>
                <div className="col-sm-12 col-xs-12 center-black">
                  <span className="btn-black" data-filter="contractNumber">No. de Contrato</span>
                </div>
              </div>

              <div className="input-group mobile-contracts">
                <input type="text" className="form-control" id="query" name="query" placeholder="Buscar" />
                <button className="btn btn-primary btn-mobile" type="submit">Buscar</button>
              </div>
            </form>
          </div>
        </div>
      </div>
    );
  }
}

export default SearchBar;
