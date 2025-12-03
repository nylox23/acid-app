import { createSlice, PayloadAction } from "@reduxjs/toolkit";
import { useSelector } from "react-redux";
import type { Acid } from "../modules/types";
import { RootState } from "../store"; // Мы создадим store на следующем шаге

interface AcidsState {
    Data: Acid[];
    SearchTerm: string;
}

const initialState: AcidsState = {
    Data: [],
    SearchTerm: ''
};

const acidsSlice = createSlice({
    name: "acids",
    initialState,
    reducers: {
        setAcids(state, action: PayloadAction<Acid[]>) {
            state.Data = action.payload;
        },
        setSearchTerm(state, action: PayloadAction<string>) {
            state.SearchTerm = action.payload;
        }
    }
});

// Селекторы (Hooks) для получения данных из store
export const useAcidsList = () =>
    useSelector((state: RootState) => state.acidsData.Data);

export const useSearchTerm = () =>
    useSelector((state: RootState) => state.acidsData.SearchTerm);

// Экспорт действий
export const {
    setAcids: setAcidsAction,
    setSearchTerm: setSearchTermAction
} = acidsSlice.actions;

export default acidsSlice.reducer;