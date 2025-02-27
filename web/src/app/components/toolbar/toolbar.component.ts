import {
  Component,
  computed,
  inject,
  OnDestroy,
  signal,
  effect,
} from '@angular/core';
import { MatToolbar } from '@angular/material/toolbar';
import { IconButtonComponent } from '../icon-button';
import { SidebarService } from '@app/services/sidebar';
import { LanguageSelectorComponent } from '../language-selector/language-selector.component';
import { TranslateModule, TranslateService } from '@ngx-translate/core';
import { BreadcrumbComponent } from '../breadcrumb';
import { BreadcrumbService } from '@app/services/breadcrumb';
import { Breadcrumb } from '@app/models';
import { Subject, takeUntil } from 'rxjs';
import { LanguageService } from '@app/services/language';

@Component({
  selector: 'app-toolbar',
  templateUrl: './toolbar.component.html',
  styleUrl: './toolbar.component.scss',
  imports: [
    MatToolbar,
    TranslateModule,
    IconButtonComponent,
    LanguageSelectorComponent,
    BreadcrumbComponent,
  ],
})
export class ToolbarComponent implements OnDestroy {
  private sidebarService = inject(SidebarService);
  private translateService = inject(TranslateService);
  private languageService = inject(LanguageService);

  isCollapsed = true;
  language = signal(this.languageService.language);

  label = computed(() => this.language());

  private breadcrumbService = inject(BreadcrumbService);
  breadcrumb: Breadcrumb[] = [];

  private onDestroy$ = new Subject();

  constructor() {
    this.breadcrumbService.breadcrumb$
      .pipe(takeUntil(this.onDestroy$))
      .subscribe(breadcrumb => {
        this.breadcrumb = breadcrumb;
      });

    effect(() => {
      this.translateService.use(this.language());
      this.languageService.update(this.language());
    });
  }

  ngOnDestroy(): void {
    this.onDestroy$.complete();
    this.onDestroy$.unsubscribe();
  }

  toggleSadebar() {
    this.sidebarService.toggleSidebar();
    this.isCollapsed = !this.isCollapsed;
  }
}
