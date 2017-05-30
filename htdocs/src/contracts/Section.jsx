import React from 'react';
import Description from '../base/Description.jsx';
import SearchResults from './SearchResults.jsx';
import SearchBar from './SearchBar.jsx';
import Details from './Details.jsx';
import Widget from './Widget.jsx';

class Section extends React.Component {
  constructor(props) {
    super(props);
    this.runQuery = this.runQuery.bind(this);
    this.state = {
      items: [],
      details: null
    };
  }

  // Submit query and update component state with results
  runQuery(q) {
    $.ajax({
      type: "POST",
      url: "/query",
      data: {
        query: JSON.stringify(q.query),
        limit: q.limit
      },
      success: (res) => this.setState({ items: res })
    });
  }

  showDetails( c ) {
    this.setState({details:c});
  }

  render() {
    if( this.props.params.id && !this.state.details) {
      $.ajax({
        type: "GET",
        url: "/contract/" + this.props.params.id ,
        success: (res) => {
          if( Object.keys(res).length > 0 ) {
            this.showDetails(res);
          }
        }
      });
    }

    if( !this.state.details ) {
      return (
        <div>
          <Description
            title="Contratos"
            color="blue aquablue"
            content="Consulta cada contrato que está registrado en Testigo Social 2.0. Podrás encontrar información para cada una de las etapas del procedimiento de contratación, desde su planeación hasta su implementación." />
          <SearchBar onSubmit={this.runQuery} />
          <SearchResults items={this.state.items} onSelection={(c) => this.showDetails(c)} />
          <div className="inner-row">
            <div className="row">
              <Widget
                title="Los 10 contratos más recientes"
                description="Ultimos contratos registrados"
                onSelection={(c) => this.showDetails(c)}
                query={{
                  query: JSON.stringify({}),
                  limit: 10,
                  sort: JSON.stringify(["-releases.date"])
                }}/>
              <Widget
                title="Los 10 contratos de mayor valor"
                description="Contratos adjudicados de mayor valor"
                onSelection={(c) => this.showDetails(c)}
                query={{
                  query: JSON.stringify({}),
                  limit: 10,
                  sort: JSON.stringify(["-releases.planning.budget.amount.amount"])
                }}/>
            </div>
          </div>
        </div>
      );
    } else {
      return (
        <div>
          <Description
            title="Contratos"
            color="blue"
            content="En esta sección podrás visualizar la información completa y agregada de un contrato, desde su planeación hasta su implementación. Como base se utiliza el Estándar de Datos para las Contrataciones Abiertas." />
          <Details contract={this.state.details} onClose={() => this.setState({details:null})} />
        </div>
      );
    }
  }
}

export default Section;
