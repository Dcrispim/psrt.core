import { useEditor } from '../context/useEditor';

export function FontsConstsSection() {
  const { state, addFont, removeFont, addConst, removeConst } = useEditor();

  return (
    <>
      <h3>Fonts</h3>
      <ul className="mini-list">
        {(state?.fonts ?? []).map((url) => (
          <li key={url}>
            <span>{url.length > 48 ? `${url.slice(0, 48)}…` : url}</span>
            <button type="button" onClick={() => removeFont(url)}>
              ×
            </button>
          </li>
        ))}
      </ul>
      <button
        type="button"
        onClick={() => {
          const url = prompt('Font URL');
          if (url) addFont(url);
        }}
      >
        Add font URL
      </button>

      <h3>Consts</h3>
      <div className="consts-table-wrap">
        <table className="consts-table">
          <thead>
            <tr>
              <th>Nome</th>
              <th>Valor</th>
              <th>Ação</th>
            </tr>
          </thead>
          <tbody>
            {Object.entries(state?.consts ?? {}).map(([name, val]) => (
              <tr key={name}>
                <td>{name}</td>
                <td>
                  <input type="text" value={val} readOnly />
                </td>
                <td>
                  <button
                    type="button"
                    onClick={() => removeConst(name)}
                  >
                    Remove
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      <button
        type="button"
        onClick={() => {
          const name = prompt('Const name');
          const value = prompt('Value');
          if (name && value != null) addConst(name, value);
        }}
      >
        Add const
      </button>
    </>
  );
}
