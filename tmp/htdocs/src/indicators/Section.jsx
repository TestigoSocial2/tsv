import React from 'react';
import { Description } from '../general.jsx';
import ChartWidget from './ChartWidget.jsx';
import DataInspector from './DataInspector.jsx';

function procedureTypeData( data ) {
  // Calculate chart total and %
  let total = data.limited.count + data.open.count + data.selective.count;
  let directP = ((data.limited.count * 100) / total).toFixed(2);
  let limitedP = ((data.selective.count * 100) / total).toFixed(2);
  let publicP = ((data.open.count * 100) / total).toFixed(2);

  // Assamble chart data
  let chartData = {
    labels: [
      'Adjudicación Directa (%)',
      'Invitación a cuando menos 3 personas (%)',
      'Licitación Pública (%)'
    ],
    datasets: [
      {
        data: [directP, limitedP, publicP],
        backgroundColor: ["#CCB3FF", "#FF6384", "#EEEEEE"],
      }
    ]
  }

  return chartData;
}

function publishYearData( data ) {
  console.log( data );
  let i = 0;
  let chartData = {
    labels: [],
    datasets: [{
      data: [],
      backgroundColor: []
    }]
  };
  let colors = ["#CCB3FF", "#FF6384", "#EEEEEE"];
  for( var y in data.years ) {
    chartData.labels.push(y);
    chartData.datasets[0].data.push( data.years[y].count );
    chartData.datasets[0].backgroundColor.push(colors[i])
    i++;
  }

  return chartData;
}

class Section extends React.Component {
  constructor(props) {
    super(props);
    this.applyFilters = this.applyFilters.bind(this);
    this.defaultFilters = {
      bucket: "gacm",
      state: "planning",
      amount: [20000000,80000000]
    };
    this.state = {
      data: {}
    };
  }

  componentDidMount() {
    this.applyFilters(this.defaultFilters);
  }

  applyFilters(filters) {
    this.runQuery(this.defaultFilters);
  }

  runQuery(filters) {
    // success: (res) => this.setState({ items: JSON.parse( res ) })
    $.ajax({
      type: "POST",
      url: "/indicators",
      data: {
        query: JSON.stringify(filters)
      },
      success: (res) => this.setState({data: JSON.parse( res )})
    });
  }

  render() {
    return (
      <div>
        <Description
          title="Indicadores"
          color="green"
          content="Testigo Social 2.0 te puede hacer llegar datos e información específica sobre procedimientos contratación pública que están en marcha. Desde un aviso de inicio de un nuevo procedimiento hasta la liga para consultar un contrato. Completa la información correspondiente y comienza a recibir notificaciones." />

        {/* Content */}
        <div className="inner-row">
          {/* Widgets */}
          <div className="row">
            <div className="col-md-6">
              <ChartWidget
                id="procedureType"
                title="Tipo de Procedimiento"
                data={this.state.data}
                reducer={procedureTypeData}
                width="500"
                height="340"
                description="La gráfica muestra la relación de contratos que se adjudicarón de acuerdo a los distintos mecanismos establecidos." />
            </div>
            <div className="col-md-6">
              <ChartWidget
                id="publishYear"
                title="Año de Publicación"
                data={this.state.data}
                reducer={publishYearData}
                width="500"
                height="340"
                description="La gráfica muestra la relación de los contratos registrados de acuerdo a su año de publicación." />
            </div>
          </div>

          {/* Separator */}
          <hr />

          {/* Main data inspector */}
          <DataInspector defaultFilters={this.defaultFilters} onSubmit={this.applyFilters}  />
        </div>
      </div>
    );
  }
}

export default Section;
