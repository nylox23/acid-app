import { useEffect } from 'react';
import { Container, Table, Spinner, Badge } from 'react-bootstrap';
import { useDispatch, useSelector } from 'react-redux';
import type { AppDispatch, RootState } from '../store';
import { fetchCarbonatesList } from '../slices/carbonateSlice';
import { useNavigate } from 'react-router-dom';
import { ROUTES } from '../Routes';
import { BreadCrumbs } from '../components/BreadCrumbs';

export const CarbonatesListPage = () => {
    const dispatch = useDispatch<AppDispatch>();
    const navigate = useNavigate();
    const { list, loading } = useSelector((state: RootState) => state.carbonateData);

    useEffect(() => {
        dispatch(fetchCarbonatesList());
    }, [dispatch]);

    const getStatusVariant = (status: string | undefined) => {
        switch (status) {
            case 'Draft': return 'warning';
            case 'Created': return 'primary';
            case 'Finished': return 'success';
            case 'Rejected': return 'danger';
            default: return 'secondary';
        }
    };

    return (
        <Container>
            <BreadCrumbs crumbs={[{ label: "Мои заявки" }]} />

            <h2 className="mb-4">История расчетов</h2>

            {loading ? (
                <div className="d-flex justify-content-center p-5">
                    <Spinner animation="border" style={{ color: '#004976' }} />
                </div>
            ) : (
                <div className="card shadow-sm">
                    <Table hover responsive className="mb-0">
                        <thead className="table-light">
                        <tr>
                            <th style={{ width: '10%' }} className="text-center">#</th>
                            <th style={{ width: '30%' }}>Дата создания</th>
                            <th style={{ width: '30%' }}>Масса CaCO3</th>
                            <th style={{ width: '30%' }} className="text-end pe-4">Статус</th>
                        </tr>
                        </thead>
                        <tbody>
                        {list.length === 0 ? (
                            <tr>
                                <td colSpan={4} className="text-center py-4 text-muted">
                                    Заявок пока нет
                                </td>
                            </tr>
                        ) : (
                            list.map(c => (
                                <tr
                                    key={c.id}
                                    onClick={() => navigate(`${ROUTES.CARBONATE_DETAIL}/${c.id}`)}
                                    style={{ cursor: 'pointer' }}
                                    className="align-middle"
                                >
                                    <td className="text-center fw-bold text-secondary">{c.id}</td>
                                    <td>
                                        {c.date_create
                                            ? new Date(c.date_create).toLocaleDateString() + ' ' + new Date(c.date_create).toLocaleTimeString([], {hour: '2-digit', minute:'2-digit'})
                                            : '-'}
                                    </td>
                                    <td>
                                        {c.mass ? `${c.mass} г` : <span className="text-muted fst-italic">Не указана</span>}
                                    </td>
                                    <td className="text-end pe-4">
                                        <Badge bg={getStatusVariant(c.status)}>
                                            {c.status === 'Draft' ? 'Черновик' : c.status}
                                        </Badge>
                                    </td>
                                </tr>
                            ))
                        )}
                        </tbody>
                    </Table>
                </div>
            )}
        </Container>
    );
};