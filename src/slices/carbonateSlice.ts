import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import { api } from '../api';
import type {
    DtoCarbonateDetailResponse,
    DtoCarbonateListEntry
} from '../api/Api';

interface CarbonateState {
    currentCarbonateId: number | null;
    currentAcidCount: number;
    detail: DtoCarbonateDetailResponse | null;
    list: DtoCarbonateListEntry[];
    loading: boolean;
    error: string | null;
}

const initialState: CarbonateState = {
    currentCarbonateId: null,
    currentAcidCount: 0,
    detail: null,
    list: [],
    loading: false,
    error: null,
};

export const fetchCurrentCarbonateInfo = createAsyncThunk(
    'carbonates/fetchCurrentInfo',
    async () => {
        // FIX: api.api.carbonatesCurrentList
        const response = await api.api.carbonatesCurrentList();
        return response.data;
    }
);

export const fetchCarbonatesList = createAsyncThunk(
    'carbonates/fetchList',
    async () => {
        // FIX: api.api.carbonatesList
        const response = await api.api.carbonatesList();
        return response.data.carbonates || [];
    }
);

export const fetchCarbonateDetail = createAsyncThunk(
    'carbonates/fetchDetail',
    async (id: number) => {
        // FIX: api.api.carbonatesDetail
        const response = await api.api.carbonatesDetail(id);
        return response.data;
    }
);

export const addAcidToCarbonate = createAsyncThunk(
    'carbonates/addAcid',
    async (acidId: number, { dispatch }) => {
        // FIX: api.api.acidsToCarbonateCreate
        await api.api.acidsToCarbonateCreate(acidId);
        dispatch(fetchCurrentCarbonateInfo());
    }
);

export const removeAcidFromCarbonate = createAsyncThunk(
    'carbonates/removeAcid',
    async (acidId: number, { dispatch, getState }) => {
        // FIX: api.api.carbonateAcidsDelete
        await api.api.carbonateAcidsDelete(acidId);

        const state = getState() as any;
        if (state.carbonateData.detail?.id) {
            dispatch(fetchCarbonateDetail(state.carbonateData.detail.id));
        }
    }
);

export const updateAcidMass = createAsyncThunk(
    'carbonates/updateAcidMass',
    async ({ id, mass }: { id: number, mass: number }, { dispatch, getState }) => {
        // FIX: api.api.carbonateAcidsUpdate
        await api.api.carbonateAcidsUpdate(id, { mass });

        const state = getState() as any;
        if (state.carbonateData.detail?.id) {
            dispatch(fetchCarbonateDetail(state.carbonateData.detail.id));
        }
    }
);

export const updateCarbonateMass = createAsyncThunk(
    'carbonates/updateCarbonateMass',
    async (mass: number, { dispatch, getState }) => {
        // FIX: api.api.carbonatesUpdate
        await api.api.carbonatesUpdate({ mass });

        const state = getState() as any;
        if (state.carbonateData.detail?.id) {
            dispatch(fetchCarbonateDetail(state.carbonateData.detail.id));
        }
    }
);

export const submitCarbonateForm = createAsyncThunk(
    'carbonates/submitForm',
    async (_, { dispatch }) => {
        // FIX: api.api.carbonatesFormUpdate
        await api.api.carbonatesFormUpdate();
        dispatch(fetchCurrentCarbonateInfo());
    }
);

export const deleteCarbonate = createAsyncThunk(
    'carbonates/delete',
    async (id: number, { dispatch }) => {
        // FIX: api.api.carbonatesDelete
        await api.api.carbonatesDelete(id);
        dispatch(fetchCurrentCarbonateInfo());
    }
);

const carbonateSlice = createSlice({
    name: 'carbonates',
    initialState,
    reducers: {
        resetDetail: (state) => { state.detail = null; }
    },
    extraReducers: (builder) => {
        builder
            .addCase(fetchCurrentCarbonateInfo.fulfilled, (state, action) => {
                state.currentCarbonateId = action.payload.carbonate_id || null;
                state.currentAcidCount = action.payload.acid_count || 0;
            })
            .addCase(fetchCarbonatesList.pending, (state) => { state.loading = true; })
            .addCase(fetchCarbonatesList.fulfilled, (state, action) => {
                state.loading = false;
                state.list = action.payload;
            })
            .addCase(fetchCarbonateDetail.pending, (state) => { state.loading = true; })
            .addCase(fetchCarbonateDetail.fulfilled, (state, action) => {
                state.loading = false;
                state.detail = action.payload;
            })
            .addCase(submitCarbonateForm.fulfilled, (state) => {
                state.currentCarbonateId = null;
                state.currentAcidCount = 0;
                state.detail = null;
            })
            .addMatcher(
                (action) => action.type === 'auth/logout/fulfilled',
                (state) => {
                    state.currentCarbonateId = null;
                    state.currentAcidCount = 0;
                    state.detail = null;
                    state.list = [];
                }
            );
    }
});

export const { resetDetail } = carbonateSlice.actions;
export default carbonateSlice.reducer;