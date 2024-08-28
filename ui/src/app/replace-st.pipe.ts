import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'replaceST',
  standalone: true,
})
export class ReplaceSTPipe implements PipeTransform {
  transform(value: string): string {
    value = value.replace(/STDOUT:\s*/g, '\n').trim();
    value = value.replace(/STDERR:\s*/g, '\n').trim();
    return value;
  }
}
