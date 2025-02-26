import { Directive, EventEmitter, HostListener, Output } from '@angular/core';

@Directive({
  selector: '[appInfiniteScroll]',
})
export class InfiniteScrollDirective {
  @Output() scrolled = new EventEmitter<void>();

  private lastScrollTop = 0;

  @HostListener('scroll', ['$event'])
  onScroll(event: Event): void {
    const target = event.target as HTMLElement;
    const isUp = target.scrollTop > this.lastScrollTop;
    if (
      isUp &&
      target.scrollTop + target.clientHeight >= target.scrollHeight - 50
    ) {
      this.scrolled.emit();
    }

    this.lastScrollTop = target.scrollTop;
  }
}
