import { Component, inject } from '@angular/core';
import { MatToolbar } from '@angular/material/toolbar';
import { IconButtonComponent } from '../icon-button';
import { SidebarService } from '@app/services/sidebar';
import { LanguageSelectorComponent } from '../language-selector/language-selector.component';
import { TranslateModule, TranslateService } from '@ngx-translate/core';

@Component({
  selector: 'app-toolbar',
  templateUrl: './toolbar.component.html',
  styleUrl: './toolbar.component.scss',
  imports: [
    MatToolbar,
    TranslateModule,
    IconButtonComponent,
    LanguageSelectorComponent,
  ],
})
export class ToolbarComponent {
  isCollapsed = true;

  private sidebarService = inject(SidebarService);
  private translateService = inject(TranslateService);

  toggleSadebar() {
    this.sidebarService.toggleSidebar();
    this.isCollapsed = !this.isCollapsed;
  }

  selectedLanguage(language: string) {
    this.translateService.use(language);
  }
}
