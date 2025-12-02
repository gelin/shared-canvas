export class PaletteChangeEvent extends Event {
    readonly #color: string;
    readonly #size: number;

    constructor(color: string, size: number) {
        super("palettechange");
        this.#color = color;
        this.#size = size;
    }

    get color() {
        return this.#color;
    }

    get size() {
        return this.#size;
    }
}
